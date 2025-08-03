package packages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"opm/db"
	"opm/middleware"
	"opm/models"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// List returns a list of packages with filtering and pagination
func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	filter := models.PackageFilter{
		Limit:  50,
		Offset: 0,
	}

	// Type filter
	if t := r.URL.Query().Get("type"); t != "" {
		pkgType := models.PackageType(t)
		filter.Type = &pkgType
	}

	// Status filter
	if s := r.URL.Query().Get("status"); s != "" {
		status := models.PackageStatus(s)
		filter.Status = &status
	}

	// Tags filter (multiple)
	if tags := r.URL.Query()["tag"]; len(tags) > 0 {
		filter.Tags = tags
	}

	// Pagination
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			filter.Limit = l
		}
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	// Build query
	query := `
			SELECT DISTINCT p.id, p.slug, p.display_name, p.description, p.type, p.status,
			       p.repository_url, p.license, p.author_id, p.created_at, p.updated_at,
			       p.view_count, p.bookmark_count,
			       u.username, u.alias, u.display_name, u.avatar_url,
			       (SELECT COUNT(*) FROM flags WHERE package_id = p.id AND status = 'pending') as active_reports_count
			FROM packages p
			JOIN users u ON p.author_id = u.id
			WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	// Apply filters
	if filter.Type != nil {
		query += fmt.Sprintf(" AND p.type = $%d", argIndex)
		args = append(args, *filter.Type)
		argIndex++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND p.status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.AuthorID != nil {
		query += fmt.Sprintf(" AND p.author_id = $%d", argIndex)
		args = append(args, *filter.AuthorID)
		argIndex++
	}

	// Tag filter - must have all specified tags
	if len(filter.Tags) > 0 {
		tagPlaceholders := []string{}
		for _, tag := range filter.Tags {
			tagPlaceholders = append(tagPlaceholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, tag)
			argIndex++
		}
		query += fmt.Sprintf(`
				AND p.id IN (
					SELECT pt.package_id 
					FROM package_tags pt
					JOIN tags t ON pt.tag_id = t.id
					WHERE t.name IN (%s)
					GROUP BY pt.package_id
					HAVING COUNT(DISTINCT t.name) = %d
				)`, strings.Join(tagPlaceholders, ","), len(filter.Tags))
	}

	// Order and pagination
	query += " ORDER BY p.created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Execute query
	rows, err := db.Conn.Query(ctx, query, args...)
	if err != nil {
		fmt.Printf("Package query error: %v\n", err)
		fmt.Printf("Query: %s\n", query)
		fmt.Printf("Args: %v\n", args)
		http.Error(w, "Failed to fetch packages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	packages := []models.Package{}
	for rows.Next() {
		var p models.Package
		var author models.User

		var activeReportsCount int
		err := rows.Scan(
			&p.ID, &p.Slug, &p.DisplayName, &p.Description, &p.Type, &p.Status,
			&p.RepositoryURL, &p.License, &p.AuthorID, &p.CreatedAt, &p.UpdatedAt,
			&p.ViewCount, &p.BookmarkCount,
			&author.Username, &author.Alias, &author.DisplayName, &author.AvatarURL,
			&activeReportsCount,
		)
		p.ActiveReportsCount = activeReportsCount
		if err != nil {
			http.Error(w, "Failed to scan package", http.StatusInternalServerError)
			return
		}
		p.Author = &author
		packages = append(packages, p)
	}

	// Get tags for each package
	for i := range packages {
		tags, err := getPackageTags(ctx, packages[i].ID, 0) // 0 for no user context
		if err == nil {
			packages[i].Tags = tags
		}
	}

	// Check bookmarks if user is authenticated
	if authUser, ok := middleware.GetAuthUser(ctx); authUser != nil {
		if ok {
			for i := range packages {
				packages[i].IsBookmarked = checkBookmark(ctx, authUser.UserID, packages[i].ID)
			}
		}

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packages)
}

// Get returns a single package by alias and slug
func Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	alias := vars["alias"]
	slug := vars["slug"]

	query := `
			SELECT p.id, p.slug, p.display_name, p.description, p.type, p.status,
			       p.repository_url, p.license, p.author_id, p.created_at, p.updated_at,
			       p.view_count, p.bookmark_count,
			       u.id, u.username, u.alias, u.display_name, u.avatar_url,
			       u.discord_verified, u.github_verified
			FROM packages p
			JOIN users u ON p.author_id = u.id
			WHERE u.alias = $1 AND p.slug = $2`

	var p models.Package
	var author models.User
	err := db.Conn.QueryRow(ctx, query, alias, slug).Scan(
		&p.ID, &p.Slug, &p.DisplayName, &p.Description, &p.Type, &p.Status,
		&p.RepositoryURL, &p.License, &p.AuthorID, &p.CreatedAt, &p.UpdatedAt,
		&p.ViewCount, &p.BookmarkCount,
		&author.ID, &author.Username, &author.Alias, &author.DisplayName, &author.AvatarURL,
		&author.DiscordVerified, &author.GitHubVerified,
	)
	if err == pgx.ErrNoRows {
		fmt.Println("ERROR-NO-ROWS", err)
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Println("GET ERR", err)
		http.Error(w, "Failed to fetch package", http.StatusInternalServerError)
		return
	}

	p.Author = &author

	// Get tags with user votes if authenticated
	userID := 0
	if authUser, ok := middleware.GetAuthUser(ctx); ok {
		userID = authUser.UserID
		p.IsBookmarked = checkBookmark(ctx, userID, p.ID)
		// fmt.Printf("Bookmark check for user %d, package %d: %v\n", userID, p.ID, p.IsBookmarked)
	} else {
		// fmt.Println("No authenticated user found in context")
	}

	tags, err := getPackageTags(ctx, p.ID, userID)
	if err == nil {
		p.Tags = tags
	}

	// Track view after successfully loading the package
	var userIDPtr *int
	if userID > 0 {
		userIDPtr = &userID
	}
	go trackView(p.ID, userIDPtr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Create creates a new package
func Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.CreatePackageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input (basic validation, could use a validation library)
	if input.Slug == "" || input.DisplayName == "" || input.Description == "" {
		http.Error(w, "Slug, display name, and description are required", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	// Check if package slug already exists
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM packages WHERE slug = $1)", input.Slug).Scan(&exists)
	if err != nil {
		http.Error(w, "Failed to check package existence", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Package slug already exists", http.StatusConflict)
		return
	}

	// Create package
	var packageID int
	err = tx.QueryRow(ctx, `
			INSERT INTO packages (slug, display_name, description, type, status, repository_url, license, author_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`,
		input.Slug, input.DisplayName, input.Description, input.Type, input.Status,
		input.RepositoryURL, input.License, authUser.UserID,
	).Scan(&packageID)
	if err != nil {
		http.Error(w, "Failed to create package", http.StatusInternalServerError)
		return
	}

	// Add tags if provided
	if len(input.TagIDs) > 0 {
		for _, tagID := range input.TagIDs {
			_, err = tx.Exec(ctx, `
					INSERT INTO package_tags (package_id, tag_id)
					VALUES ($1, $2)
					ON CONFLICT DO NOTHING`,
				packageID, tagID,
			)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to add tags: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Return the created package
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   packageID,
		"slug": input.Slug,
	})
}

// Helper functions

// prettyPrint formats any struct for debug logging
func prettyPrint(label string, v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("%s: %+v (error formatting: %v)\n", label, v, err)
		return
	}
	fmt.Printf("%s:\n%s\n", label, string(b))
}

func getPackageTags(ctx context.Context, packageID int, userID int) ([]models.Tag, error) {
	query := `
		SELECT t.id, t.name, t.usage_count, t.created_at, 
		       pt.score as net_score,
		       COALESCE(tv.vote_value, 0) as user_vote
		FROM tags t
		JOIN package_tags pt ON t.id = pt.tag_id
		LEFT JOIN tag_votes tv ON pt.package_id = tv.package_id 
		                       AND pt.tag_id = tv.tag_id 
		                       AND tv.user_id = $2
		WHERE pt.package_id = $1
		ORDER BY pt.score DESC, t.name`

	rows, err := db.Conn.Query(ctx, query, packageID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		var t models.Tag
		err := rows.Scan(&t.ID, &t.Name, &t.UsageCount, &t.CreatedAt, &t.NetScore, &t.UserVote)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func checkBookmark(ctx context.Context, userID, packageID int) bool {
	var exists bool
	err := db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND package_id = $2)",
		userID, packageID,
	).Scan(&exists)
	if err != nil {
		fmt.Printf("Error checking bookmark: %v\n", err)
	}
	return err == nil && exists
}

func trackView(packageID int, userID *int) {
	// Use a new context with timeout to not block the request
	trackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call the stored procedure
	_, err := db.Conn.Exec(trackCtx, "SELECT track_package_view($1, $2)", packageID, userID)
	if err != nil {
		// Log but don't fail
		fmt.Printf("Failed to track view for package %d: %v\n", packageID, err)
	}
}

// Update updates a package
func Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	packageID := vars["id"]

	// Verify ownership
	var authorID int
	err := db.Conn.QueryRow(ctx,
		"SELECT author_id FROM packages WHERE id = $1",
		packageID,
	).Scan(&authorID)
	if err == pgx.ErrNoRows {
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to find package", http.StatusInternalServerError)
		return
	}

	if authorID != authUser.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse update input
	var input models.UpdatePackageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Debug log
	prettyPrint("Update input", input)

	// Build update query dynamically
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if input.DisplayName != nil {
		updateFields = append(updateFields, fmt.Sprintf("display_name = $%d", argIndex))
		args = append(args, *input.DisplayName)
		argIndex++
	}

	if input.Description != nil {
		updateFields = append(updateFields, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *input.Description)
		argIndex++
	}

	if input.Status != nil {
		updateFields = append(updateFields, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *input.Status)
		argIndex++
	}

	if input.RepositoryURL != nil {
		updateFields = append(updateFields, fmt.Sprintf("repository_url = $%d", argIndex))
		args = append(args, *input.RepositoryURL)
		argIndex++
	}

	if input.License != nil {
		if *input.License == "" {
			// Empty string means clear the license (set to NULL)
			updateFields = append(updateFields, fmt.Sprintf("license = NULL"))
		} else {
			updateFields = append(updateFields, fmt.Sprintf("license = $%d", argIndex))
			args = append(args, *input.License)
			argIndex++
		}
	}

	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Add package ID as last argument
	args = append(args, packageID)

	query := fmt.Sprintf(
		"UPDATE packages SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = $%d",
		strings.Join(updateFields, ", "),
		argIndex,
	)

	// Debug log the query
	fmt.Printf("Update query: %s\n", query)
	prettyPrint("Update args", args)

	_, err = db.Conn.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("Update error: %v\n", err)
		http.Error(w, "Failed to update package", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// Delete deletes a package
func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	packageID := vars["id"]

	// Verify ownership
	var authorID int
	err := db.Conn.QueryRow(ctx,
		"SELECT author_id FROM packages WHERE id = $1",
		packageID,
	).Scan(&authorID)
	if err == pgx.ErrNoRows {
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to find package", http.StatusInternalServerError)
		return
	}

	if authorID != authUser.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete package (cascades to related tables)
	_, err = db.Conn.Exec(ctx, "DELETE FROM packages WHERE id = $1", packageID)
	if err != nil {
		http.Error(w, "Failed to delete package", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
