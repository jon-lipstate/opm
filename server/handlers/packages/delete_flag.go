package packages

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"opm/db"
	"opm/middleware"
)

// DeleteFlag allows a user to delete their own flag
func DeleteFlag(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to delete flag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}