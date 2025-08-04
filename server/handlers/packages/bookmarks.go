package packages

import (
	"encoding/json"
	"net/http"
	"opm/db"
	"opm/helpers"
	"opm/logger"
	"opm/middleware"
)

// Bookmark adds a bookmark for the authenticated user
func Bookmark(w http.ResponseWriter, r *http.Request) {
	mainLogger := logger.MainLogger
	ctx := r.Context()

	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	packageID, ok := helpers.RequiredParamInt(r, w, "package_id")
	if !ok {
		return
	}

	// Add bookmark
	_, err := db.Conn.Exec(ctx,
		"INSERT INTO bookmarks (user_id, package_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		authUser.UserID, packageID,
	)
	if err != nil {
		mainLogger.Printf("Failed to add bookmark for user %d, package %d: %v", authUser.UserID, packageID, err)
		http.Error(w, "Failed to add bookmark", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// Unbookmark removes a bookmark for the authenticated user
func Unbookmark(w http.ResponseWriter, r *http.Request) {
	mainLogger := logger.MainLogger
	ctx := r.Context()

	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		mainLogger.Println("not auth")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	packageID, ok := helpers.RequiredParamInt(r, w, "package_id")
	if !ok {
		mainLogger.Println("MISSING package_id")
		return
	}

	// Remove bookmark
	_, err := db.Conn.Exec(ctx,
		"DELETE FROM bookmarks WHERE user_id = $1 AND package_id = $2",
		authUser.UserID, packageID,
	)
	if err != nil {
		mainLogger.Printf("Failed to remove bookmark for user %d, package %d: %v", authUser.UserID, packageID, err)
		http.Error(w, "Failed to remove bookmark", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
