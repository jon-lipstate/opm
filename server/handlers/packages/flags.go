package packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opm/db"
	"opm/middleware"
	"opm/models"
	"time"

	"github.com/gorilla/mux"
)

// FlagPackage creates a moderation flag for a package
func FlagPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		fmt.Println("Not Auth")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		PackageID int     `json:"package_id"`
		Reason    string  `json:"reason"`
		Details   *string `json:"details,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.PackageID == 0 {
		fmt.Println("Invalid Package ID")
		http.Error(w, "Package ID is required", http.StatusBadRequest)
		return
	}

	// Validate reason
	validReasons := map[string]bool{
		"Malicious code":        true,
		"Copyright violation":   true,
		"Inappropriate content": true,
		"Broken/non-functional": true,
		"Spam":                  true,
		"Other":                 true,
	}
	if !validReasons[input.Reason] {
		fmt.Println("Invalid Reason Code")
		http.Error(w, "Invalid flag reason", http.StatusBadRequest)
		return
	}

	// Verify package exists
	var packageExists bool
	err := db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM packages WHERE id = $1)",
		input.PackageID,
	).Scan(&packageExists)
	if err != nil || !packageExists {
		fmt.Println("Package Search Error", err)
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	// Check if user already flagged this package
	var existingFlag bool
	err = db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM flags WHERE package_id = $1 AND user_id = $2 AND status = 'pending')",
		input.PackageID, authUser.UserID,
	).Scan(&existingFlag)
	if err == nil && existingFlag {
		fmt.Println("Package Flags Error", err)
		http.Error(w, "You have already flagged this package", http.StatusConflict)
		return
	}

	// Create flag
	var flagID int
	err = db.Conn.QueryRow(ctx, `
			INSERT INTO flags (package_id, user_id, reason, details, status)
			VALUES ($1, $2, $3, $4, 'pending')
			RETURNING id`,
		input.PackageID, authUser.UserID, input.Reason, input.Details,
	).Scan(&flagID)
	if err != nil {
		fmt.Println("Insert Flags Error", err)
		http.Error(w, "Failed to create flag", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     flagID,
		"status": "pending",
	})
}

// GetPackageFlags returns active flags for a package (public)
func GetPackageFlags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get package ID from query parameter
	packageIDStr := r.URL.Query().Get("package_id")
	if packageIDStr == "" {
		http.Error(w, "Missing package_id parameter", http.StatusBadRequest)
		return
	}

	var packageID int
	if _, err := fmt.Sscanf(packageIDStr, "%d", &packageID); err != nil {
		http.Error(w, "Invalid package_id", http.StatusBadRequest)
		return
	}

	// Verify package exists
	var exists bool
	err := db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM packages WHERE id = $1)",
		packageID,
	).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	// Get active flags
	query := `
			SELECT id, package_id, user_id, reason, details, status, 
			       resolved_by, resolved_at, created_at, updated_at
			FROM flags
			WHERE package_id = $1 AND status = 'pending'
			ORDER BY created_at DESC`

	fmt.Printf("Fetching flags for package %d\n", packageID)
	rows, err := db.Conn.Query(ctx, query, packageID)
	if err != nil {
		fmt.Printf("Error fetching flags: %v\n", err)
		http.Error(w, "Failed to fetch flags", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	flags := []models.Flag{}
	_, ok := middleware.GetAuthUser(ctx)

	for rows.Next() {
		var f models.Flag
		err := rows.Scan(&f.ID, &f.PackageID, &f.UserID, &f.Reason, &f.Details, &f.Status, 
			&f.ResolvedBy, &f.ResolvedAt, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			fmt.Printf("Error scanning flag: %v\n", err)
			continue
		}

		// Only show details to authenticated users
		if !ok {
			f.Details = nil
		}

		flags = append(flags, f)
	}

	fmt.Printf("Found %d flags for package %d\n", len(flags), packageID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flags)
}

// GetAllFlags returns all flags (moderator only)
func GetAllFlags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is moderator
	var isModerator bool
	err := db.Conn.QueryRow(ctx,
		"SELECT is_moderator FROM users WHERE id = $1",
		authUser.UserID,
	).Scan(&isModerator)
	if err != nil || !isModerator {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get filter parameters
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "pending"
	}

	// Build query
	query := `
		SELECT f.id, f.package_id, f.user_id, f.reason, f.details, f.status, 
		       f.resolved_by, f.resolved_at, f.created_at,
		       p.slug, p.display_name,
		       u.username as reporter_username,
		       ru.username as resolver_username
		FROM flags f
		JOIN packages p ON f.package_id = p.id
		JOIN users u ON f.user_id = u.id
		LEFT JOIN users ru ON f.resolved_by = ru.id
		WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if status != "all" {
		query += fmt.Sprintf(" AND f.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	query += " ORDER BY f.created_at DESC"

	rows, err := db.Conn.Query(ctx, query, args...)
	if err != nil {
		http.Error(w, "Failed to fetch flags", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type FlagWithContext struct {
		models.Flag
		PackageSlug        string  `json:"package_slug"`
		PackageDisplayName string  `json:"package_display_name"`
		ReporterUsername   string  `json:"reporter_username"`
		ResolverUsername   *string `json:"resolver_username,omitempty"`
	}

	flags := []FlagWithContext{}
	for rows.Next() {
		var f FlagWithContext
		err := rows.Scan(
			&f.ID, &f.PackageID, &f.UserID, &f.Reason, &f.Details, &f.Status,
			&f.ResolvedBy, &f.ResolvedAt, &f.CreatedAt,
			&f.PackageSlug, &f.PackageDisplayName,
			&f.ReporterUsername, &f.ResolverUsername,
		)
		if err != nil {
			continue
		}
		flags = append(flags, f)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flags)
}

// ResolveFlag updates a flag's status (moderator only)
func ResolveFlag(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is moderator
	var isModerator bool
	err := db.Conn.QueryRow(ctx,
		"SELECT is_moderator FROM users WHERE id = $1",
		authUser.UserID,
	).Scan(&isModerator)
	if err != nil || !isModerator {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	flagID := vars["id"]

	var input struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"reviewed":  true,
		"resolved":  true,
		"dismissed": true,
	}
	if !validStatuses[input.Status] {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	// Update flag
	_, err = db.Conn.Exec(ctx, `
		UPDATE flags 
		SET status = $1, resolved_by = $2, resolved_at = CURRENT_TIMESTAMP
		WHERE id = $3`,
		input.Status, authUser.UserID, flagID,
	)
	if err != nil {
		http.Error(w, "Failed to update flag", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": input.Status,
	})
}

// GetUserFlags returns flags created by the authenticated user
func GetUserFlags(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT f.id, f.package_id, f.reason, f.details, f.status, 
		       f.resolved_at, f.created_at,
		       p.slug, p.display_name, u.alias as author_alias
		FROM flags f
		JOIN packages p ON f.package_id = p.id
		JOIN users u ON p.author_id = u.id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC`

	rows, err := db.Conn.Query(ctx, query, authUser.UserID)
	if err != nil {
		http.Error(w, "Failed to fetch flags", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserFlag struct {
		ID                 int        `json:"id"`
		PackageID          int        `json:"package_id"`
		Reason             string     `json:"reason"`
		Details            *string    `json:"details,omitempty"`
		Status             string     `json:"status"`
		ResolvedAt         *time.Time `json:"resolved_at,omitempty"`
		CreatedAt          time.Time  `json:"created_at"`
		PackageSlug        string     `json:"package_slug"`
		PackageDisplayName string     `json:"package_display_name"`
		AuthorAlias        string     `json:"author_alias"`
	}

	flags := []UserFlag{}
	for rows.Next() {
		var f UserFlag
		err := rows.Scan(
			&f.ID, &f.PackageID, &f.Reason, &f.Details, &f.Status,
			&f.ResolvedAt, &f.CreatedAt,
			&f.PackageSlug, &f.PackageDisplayName, &f.AuthorAlias,
		)
		if err != nil {
			continue
		}
		flags = append(flags, f)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flags)
}

// GetFlagStats returns flag statistics for a package
func GetFlagStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get package ID from query parameter
	packageIDStr := r.URL.Query().Get("package_id")
	if packageIDStr == "" {
		http.Error(w, "Missing package_id parameter", http.StatusBadRequest)
		return
	}

	var packageID int
	if _, err := fmt.Sscanf(packageIDStr, "%d", &packageID); err != nil {
		http.Error(w, "Invalid package_id", http.StatusBadRequest)
		return
	}

	// Get flag statistics
	query := `
		SELECT 
			COUNT(*) FILTER (WHERE status = 'pending') as pending_count,
			COUNT(*) FILTER (WHERE status = 'reviewed') as reviewed_count,
			COUNT(*) FILTER (WHERE status = 'resolved') as resolved_count,
			COUNT(*) FILTER (WHERE status = 'dismissed') as dismissed_count,
			COUNT(*) as total_count,
			COUNT(DISTINCT reason) as unique_reasons,
			MAX(created_at) as last_flag_date
		FROM flags
		WHERE package_id = $1`

	var stats struct {
		PendingCount   int        `json:"pending_count"`
		ReviewedCount  int        `json:"reviewed_count"`
		ResolvedCount  int        `json:"resolved_count"`
		DismissedCount int        `json:"dismissed_count"`
		TotalCount     int        `json:"total_count"`
		UniqueReasons  int        `json:"unique_reasons"`
		LastFlagDate   *time.Time `json:"last_flag_date,omitempty"`
	}

	err := db.Conn.QueryRow(ctx, query, packageID).Scan(
		&stats.PendingCount,
		&stats.ReviewedCount,
		&stats.ResolvedCount,
		&stats.DismissedCount,
		&stats.TotalCount,
		&stats.UniqueReasons,
		&stats.LastFlagDate,
	)
	if err != nil {
		http.Error(w, "Failed to fetch flag statistics", http.StatusInternalServerError)
		return
	}

	// Get reason breakdown
	reasonQuery := `
		SELECT reason, COUNT(*) as count
		FROM flags
		WHERE package_id = $1
		GROUP BY reason
		ORDER BY count DESC`

	rows, err := db.Conn.Query(ctx, reasonQuery, packageID)
	if err == nil {
		defer rows.Close()

		reasons := []struct {
			Reason string `json:"reason"`
			Count  int    `json:"count"`
		}{}

		for rows.Next() {
			var r struct {
				Reason string `json:"reason"`
				Count  int    `json:"count"`
			}
			if err := rows.Scan(&r.Reason, &r.Count); err == nil {
				reasons = append(reasons, r)
			}
		}

		// Add reasons to response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"stats":   stats,
			"reasons": reasons,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
