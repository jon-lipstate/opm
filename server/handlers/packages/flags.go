package packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opm/db"
	"opm/helpers"
	"opm/logger"
	"opm/middleware"
	"opm/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// FlagPackage creates a moderation flag for a package
func FlagPackage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		PackageID int     `json:"package_id"`
		Reason    string  `json:"reason"`
		Details   *string `json:"details,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if input.PackageID == 0 {
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
		http.Error(w, "Invalid flag reason", http.StatusBadRequest)
		return
	}

	// Verify package exists
	var packageExists bool
	err := db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM packages WHERE id = $1)",
		input.PackageID,
	).Scan(&packageExists)
	if err != nil {
		logger.MainLogger.Printf("Failed to check package existence for package %d: %v", input.PackageID, err)
		http.Error(w, "Failed to check package existence", http.StatusInternalServerError)
		return
	}
	if !packageExists {
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}

	// Check if user already flagged this package
	var existingFlag bool
	err = db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM flags WHERE package_id = $1 AND user_id = $2 AND status = 'pending')",
		input.PackageID, authUser.UserID,
	).Scan(&existingFlag)
	if err != nil {
		logger.MainLogger.Printf("Failed to check existing flag for package %d, user %d: %v", input.PackageID, authUser.UserID, err)
		http.Error(w, "Failed to check existing flag", http.StatusInternalServerError)
		return
	}
	if existingFlag {
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
		logger.MainLogger.Printf("Failed to create flag for package %d, user %d: %v", input.PackageID, authUser.UserID, err)
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

	packageID, ok := helpers.RequiredParamInt(r, w, "package_id")
	if !ok {
		return
	}

	// Verify package exists
	var exists bool
	err := db.Conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM packages WHERE id = $1)",
		packageID,
	).Scan(&exists)
	if err != nil {
		logger.MainLogger.Printf("Failed to check package existence for package %d: %v", packageID, err)
		http.Error(w, "Failed to check package existence", http.StatusInternalServerError)
		return
	}
	if !exists {
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

	rows, err := db.Conn.Query(ctx, query, packageID)
	if err != nil {
		logger.MainLogger.Printf("Failed to fetch flags for package %d: %v", packageID, err)
		http.Error(w, "Failed to fetch flags", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	flags := []models.Flag{}
	_, ok = middleware.GetAuthUser(ctx)

	for rows.Next() {
		var f models.Flag
		err := rows.Scan(&f.ID, &f.PackageID, &f.UserID, &f.Reason, &f.Details, &f.Status,
			&f.ResolvedBy, &f.ResolvedAt, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			logger.MainLogger.Printf("Failed to scan flag: %v", err)
			continue
		}

		// Only show details to authenticated users
		if !ok {
			f.Details = nil
		}

		flags = append(flags, f)
	}

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
	if err != nil {
		logger.MainLogger.Printf("Failed to check moderator status for user %d: %v", authUser.UserID, err)
		http.Error(w, "Failed to check permissions", http.StatusInternalServerError)
		return
	}
	if !isModerator {
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
		logger.MainLogger.Printf("Failed to fetch all flags: %v", err)
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
			logger.MainLogger.Printf("Failed to scan flag with context: %v", err)
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
	if err != nil {
		logger.MainLogger.Printf("Failed to check moderator status for user %d: %v", authUser.UserID, err)
		http.Error(w, "Failed to check permissions", http.StatusInternalServerError)
		return
	}
	if !isModerator {
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
		logger.MainLogger.Printf("Failed to update flag %s: %v", flagID, err)
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
		       p.slug, p.display_name, u.slug as author_slug
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
		AuthorSlug         string     `json:"author_slug"`
	}

	flags := []UserFlag{}
	for rows.Next() {
		var f UserFlag
		err := rows.Scan(
			&f.ID, &f.PackageID, &f.Reason, &f.Details, &f.Status,
			&f.ResolvedAt, &f.CreatedAt,
			&f.PackageSlug, &f.PackageDisplayName, &f.AuthorSlug,
		)
		if err != nil {
			logger.MainLogger.Printf("Failed to scan user flag: %v", err)
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

	packageID, ok := helpers.RequiredParamInt(r, w, "package_id")
	if !ok {
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
		logger.MainLogger.Printf("Failed to fetch flag statistics for package %d: %v", packageID, err)
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
			if err := rows.Scan(&r.Reason, &r.Count); err != nil {
				logger.MainLogger.Printf("Failed to scan flag reason: %v", err)
			} else {
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

func DeleteFlag(w http.ResponseWriter, r *http.Request) {
	mainLogger := logger.MainLogger
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get flag ID from URL
	vars := mux.Vars(r)
	flagID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid flag ID", http.StatusBadRequest)
		return
	}

	// Verify the flag belongs to the user
	var userID int
	query := `SELECT user_id FROM flags WHERE id = $1`
	err = db.QueryRow(ctx, query, flagID).Scan(&userID)
	if err != nil {
		http.Error(w, "Flag not found", http.StatusNotFound)
		return
	}

	if userID != authUser.UserID {
		http.Error(w, "Forbidden - you can only delete your own flags", http.StatusForbidden)
		return
	}

	// Delete the flag
	deleteQuery := `DELETE FROM flags WHERE id = $1`
	_, err = db.Exec(ctx, deleteQuery, flagID)
	if err != nil {
		mainLogger.Println("delete flag err", err)
		http.Error(w, "Failed to delete flag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
