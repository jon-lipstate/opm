package users

import (
	"encoding/json"
	"net/http"

	"opm/db"
	"opm/logger"
	"opm/middleware"
	"opm/models"
)

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get auth user from context
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch user details from database
	var user models.User
	query := `SELECT id, github_id, discord_id, username, slug, display_name, avatar_url, created_at, updated_at 
				  FROM users WHERE id = $1`

	err := db.QueryRow(ctx, query, authUser.UserID).Scan(
		&user.ID,
		&user.GitHubID,
		&user.DiscordID,
		&user.Username,
		&user.Slug,
		&user.DisplayName,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		logger.MainLogger.Printf("Failed to fetch user %d from database: %v", authUser.UserID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
