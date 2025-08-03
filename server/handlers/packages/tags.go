package packages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"opm/db"
	"opm/middleware"
)

// AddTag adds a tag to a package
func AddTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		PackageID int    `json:"package_id"`
		TagName   string `json:"tag_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.PackageID == 0 || input.TagName == "" {
		fmt.Println("Bad Input", input)
		http.Error(w, "Package ID and tag name are required", http.StatusBadRequest)
		return
	}

	// Verify package exists
	var packageExists bool
	err := db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM packages WHERE id = $1)",
		input.PackageID,
	).Scan(&packageExists)
	if err != nil || !packageExists {
		fmt.Println("Package not found")
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	// Get or create tag
	var tagID int
	err = db.Conn.QueryRow(ctx,
		"INSERT INTO tags (name, added_by) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id",
		input.TagName, authUser.UserID,
	).Scan(&tagID)
	if err != nil {
		fmt.Println("Insert tags Error", err)
		http.Error(w, "Failed to create tag", http.StatusInternalServerError)
		return
	}

	// Add tag to package if not already present
	_, err = db.Conn.Exec(ctx,
		"INSERT INTO package_tags (package_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		input.PackageID, tagID,
	)
	if err != nil {
		fmt.Println("Insert package_tags Error", err)
		http.Error(w, "Failed to add tag to package", http.StatusInternalServerError)
		return
	}

	// All votes are worth 1 for now
	voteValue := 1

	// Add initial vote
	_, err = db.Conn.Exec(ctx, `
			INSERT INTO tag_votes (package_id, tag_id, user_id, vote_value)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (package_id, tag_id, user_id) 
			DO UPDATE SET vote_value = EXCLUDED.vote_value`,
		input.PackageID, tagID, authUser.UserID, voteValue,
	)
	if err != nil {
		fmt.Println("Insert tag_votes Error", err)
		http.Error(w, "Failed to add vote", http.StatusInternalServerError)
		return
	}

	// Update package_tags score
	updateTagScore(ctx, input.PackageID, tagID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tag_id":     tagID,
		"tag_name":   input.TagName,
		"vote_value": voteValue,
	})
}

// VoteTag votes on a tag for a package
func VoteTag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		fmt.Println("Not Auth")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		PackageID int `json:"package_id"`
		TagID     int `json:"tag_id"`
		Vote      int `json:"vote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.PackageID == 0 || input.TagID == 0 {
		fmt.Println("Invalid JSON")
		http.Error(w, "Package ID and tag ID are required", http.StatusBadRequest)
		return
	}

	// Validate vote value
	if input.Vote < -1 || input.Vote > 1 {
		fmt.Println("Invalid Payload")
		http.Error(w, "Vote must be -1, 0, or 1", http.StatusBadRequest)
		return
	}

	// Verify package exists
	var packageExists bool
	err := db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM packages WHERE id = $1)",
		input.PackageID,
	).Scan(&packageExists)
	if err != nil || !packageExists {
		fmt.Println("Missing Package")
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	// Check if tag exists on package
	var exists bool
	err = db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM package_tags WHERE package_id = $1 AND tag_id = $2)",
		input.PackageID, input.TagID,
	).Scan(&exists)
	if err != nil || !exists {
		fmt.Println("Missing Tag for Package")
		http.Error(w, "Tag not found on package", http.StatusNotFound)
		return
	}

	// All votes are worth their base value (1 or -1)
	voteValue := input.Vote

	if voteValue == 0 {
		// Remove vote
		_, err = db.Conn.Exec(ctx,
			"DELETE FROM tag_votes WHERE package_id = $1 AND tag_id = $2 AND user_id = $3",
			input.PackageID, input.TagID, authUser.UserID,
		)
	} else {
		// Add or update vote
		_, err = db.Conn.Exec(ctx, `
				INSERT INTO tag_votes (package_id, tag_id, user_id, vote_value)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (package_id, tag_id, user_id) 
				DO UPDATE SET vote_value = EXCLUDED.vote_value`,
			input.PackageID, input.TagID, authUser.UserID, voteValue,
		)
	}

	if err != nil {
		fmt.Println("Update tag_votes error", err)
		http.Error(w, "Failed to update vote", http.StatusInternalServerError)
		return
	}

	// Update package_tags score
	newScore := updateTagScore(ctx, input.PackageID, input.TagID)

	// If score is <= 0, remove the tag from the package
	if newScore <= 0 {
		removeTagFromPackage(ctx, input.PackageID, input.TagID)
		newScore = 0 // Set to 0 to indicate removal
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"vote":       input.Vote,
		"vote_value": voteValue,
		"net_score":  newScore,
		"removed":    newScore == 0,
	})
}

// Helper function to update tag score
func updateTagScore(ctx context.Context, packageID, tagID int) int {
	var score int
	err := db.Conn.QueryRow(ctx, `
		UPDATE package_tags 
		SET score = (
			SELECT COALESCE(SUM(vote_value), 0) 
			FROM tag_votes 
			WHERE package_id = $1 AND tag_id = $2
		)
		WHERE package_id = $1 AND tag_id = $2
		RETURNING score`,
		packageID, tagID,
	).Scan(&score)

	if err != nil {
		fmt.Println("Update Package_Tags score error")
		return 0
	}
	return score
}

// Helper function to remove tag from package and clean up orphaned tags
func removeTagFromPackage(ctx context.Context, packageID, tagID int) {
	// Start a transaction
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		fmt.Printf("Failed to start transaction for tag removal: %v\n", err)
		return
	}
	defer tx.Rollback(ctx)

	// Remove all votes for this tag on this package
	_, err = tx.Exec(ctx,
		"DELETE FROM tag_votes WHERE package_id = $1 AND tag_id = $2",
		packageID, tagID,
	)
	if err != nil {
		fmt.Printf("Failed to remove tag votes: %v\n", err)
		return
	}

	// Remove the tag from the package
	_, err = tx.Exec(ctx,
		"DELETE FROM package_tags WHERE package_id = $1 AND tag_id = $2",
		packageID, tagID,
	)
	if err != nil {
		fmt.Printf("Failed to remove tag from package: %v\n", err)
		return
	}

	// Check if this tag is used by any other packages
	var isUsed bool
	err = tx.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM package_tags WHERE tag_id = $1)",
		tagID,
	).Scan(&isUsed)
	if err != nil {
		fmt.Printf("Failed to check tag usage: %v\n", err)
		return
	}

	// If tag is not used by any other packages, delete it
	if !isUsed {
		_, err = tx.Exec(ctx, "DELETE FROM tags WHERE id = $1", tagID)
		if err != nil {
			fmt.Printf("Failed to delete orphaned tag: %v\n", err)
			return
		}
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		fmt.Printf("Failed to commit tag removal transaction: %v\n", err)
	}
}
