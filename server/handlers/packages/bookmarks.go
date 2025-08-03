package packages

import (
	"encoding/json"
	"net/http"
	"opm/db"
	"opm/middleware"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// Bookmark adds a bookmark for the authenticated user
func Bookmark(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authUser, ok := middleware.GetAuthUser(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		alias := vars["alias"]
		slug := vars["slug"]

		// Get package ID
		var packageID int
		err := db.Conn.QueryRow(ctx,
			"SELECT p.id FROM packages p JOIN users u ON p.author_id = u.id WHERE u.alias = $1 AND p.slug = $2",
			alias, slug,
		).Scan(&packageID)
		if err == pgx.ErrNoRows {
			http.Error(w, "Package not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Failed to find package", http.StatusInternalServerError)
			return
		}

		// Add bookmark
		_, err = db.Conn.Exec(ctx,
			"INSERT INTO bookmarks (user_id, package_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			authUser.UserID, packageID,
		)
		if err != nil {
			http.Error(w, "Failed to add bookmark", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// Unbookmark removes a bookmark for the authenticated user
func Unbookmark(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authUser, ok := middleware.GetAuthUser(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		alias := vars["alias"]
		slug := vars["slug"]

		// Get package ID
		var packageID int
		err := db.Conn.QueryRow(ctx,
			"SELECT p.id FROM packages p JOIN users u ON p.author_id = u.id WHERE u.alias = $1 AND p.slug = $2",
			alias, slug,
		).Scan(&packageID)
		if err == pgx.ErrNoRows {
			http.Error(w, "Package not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Failed to find package", http.StatusInternalServerError)
			return
		}

		// Remove bookmark
		_, err = db.Conn.Exec(ctx,
			"DELETE FROM bookmarks WHERE user_id = $1 AND package_id = $2",
			authUser.UserID, packageID,
		)
		if err != nil {
			http.Error(w, "Failed to remove bookmark", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
