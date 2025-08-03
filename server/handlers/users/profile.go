package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opm/db"
	"opm/middleware"
	"opm/models"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

var aliasRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-_]*[a-z0-9]$`)

// isValidAlias checks if an alias is URL-safe and follows our rules
func isValidAlias(alias string) bool {
	// Must be lowercase
	if alias != strings.ToLower(alias) {
		return false
	}
	// Must match pattern: start/end with alphanumeric, can contain hyphens/underscores
	return aliasRegex.MatchString(alias)
}

// UpdateProfile updates the authenticated user's profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Build update query dynamically
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	// Validate and update alias
	if input.Alias != nil {
		// Validate alias format
		alias := strings.TrimSpace(*input.Alias)
		if alias == "" {
			http.Error(w, "Alias cannot be empty", http.StatusBadRequest)
			return
		}
		if len(alias) < 3 || len(alias) > 50 {
			http.Error(w, "Alias must be between 3 and 50 characters", http.StatusBadRequest)
			return
		}
		// Check if alias matches required pattern (lowercase letters, numbers, hyphens, underscores)
		// Must be URL-safe without encoding
		if !isValidAlias(alias) {
			http.Error(w, "Alias must contain only lowercase letters, numbers, hyphens, and underscores", http.StatusBadRequest)
			return
		}

		// Check if alias is already taken by another user
		var existingUserID int
		err := db.Conn.QueryRow(ctx,
			"SELECT id FROM users WHERE alias = $1 AND id != $2",
			alias, authUser.UserID,
		).Scan(&existingUserID)
		if err != pgx.ErrNoRows {
			http.Error(w, "Alias is already taken", http.StatusConflict)
			return
		}

		updateFields = append(updateFields, fmt.Sprintf("alias = $%d", argIndex))
		args = append(args, alias)
		argIndex++
	}

	// Update display name
	if input.DisplayName != nil {
		displayName := strings.TrimSpace(*input.DisplayName)
		if len(displayName) > 255 {
			http.Error(w, "Display name must be less than 255 characters", http.StatusBadRequest)
			return
		}
		updateFields = append(updateFields, fmt.Sprintf("display_name = $%d", argIndex))
		args = append(args, displayName)
		argIndex++
	}

	// Update avatar URL
	if input.AvatarURL != nil {
		avatarURL := strings.TrimSpace(*input.AvatarURL)
		if avatarURL != "" && !strings.HasPrefix(avatarURL, "http") {
			http.Error(w, "Avatar URL must be a valid URL", http.StatusBadRequest)
			return
		}
		updateFields = append(updateFields, fmt.Sprintf("avatar_url = $%d", argIndex))
		args = append(args, avatarURL)
		argIndex++
	}

	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Add user ID as last argument
	args = append(args, authUser.UserID)

	query := fmt.Sprintf(
		"UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = $%d",
		strings.Join(updateFields, ", "),
		argIndex,
	)

	_, err := db.Conn.Exec(ctx, query, args...)
	if err != nil {
		fmt.Printf("Update profile error: %v\n", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
	})
}

// CheckAliasAvailability checks if an alias is available
func CheckAliasAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alias := r.URL.Query().Get("alias")

	if alias == "" {
		http.Error(w, "Alias parameter is required", http.StatusBadRequest)
		return
	}

	// Validate alias format
	if !isValidAlias(alias) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": false,
			"reason":    "Invalid format. Use lowercase letters, numbers, hyphens, and underscores. Must start and end with a letter or number.",
		})
		return
	}

	// Check if alias is too short or too long
	if len(alias) < 3 || len(alias) > 50 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": false,
			"reason":    "Alias must be between 3 and 50 characters",
		})
		return
	}

	// Get current user ID if authenticated
	var currentUserID int
	if authUser, ok := middleware.GetAuthUser(ctx); ok {
		currentUserID = authUser.UserID
	}

	// Check if alias exists (excluding current user)
	var existingID int
	var query string
	var args []interface{}
	
	if currentUserID > 0 {
		query = "SELECT id FROM users WHERE alias = $1 AND id != $2"
		args = []interface{}{alias, currentUserID}
	} else {
		query = "SELECT id FROM users WHERE alias = $1"
		args = []interface{}{alias}
	}
	
	err := db.Conn.QueryRow(ctx, query, args...).Scan(&existingID)
	
	if err == pgx.ErrNoRows {
		// Alias is available
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": true,
		})
		return
	}

	// Alias is taken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"available": false,
		"reason":    "This alias is already taken",
	})
}
