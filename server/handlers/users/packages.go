package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opm/db"
	"opm/middleware"
	"opm/models"
)

// ListUserPackages returns all packages created by the authenticated user
func ListUserPackages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser, ok := middleware.GetAuthUser(ctx)
	if !ok {
		fmt.Println("Not Auth")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT p.id, p.slug, p.display_name, p.description, p.type, p.status,
		       p.repository_url, p.license, p.author_id, p.created_at, p.updated_at,
		       p.view_count, p.bookmark_count,
		       u.username, u.alias, u.avatar_url
		FROM packages p
		JOIN users u ON p.author_id = u.id
		WHERE p.author_id = $1
		ORDER BY p.created_at DESC`

	rows, err := db.Conn.Query(ctx, query, authUser.UserID)
	if err != nil {
		fmt.Println("Search packages Err", err)
		http.Error(w, "Failed to fetch packages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	packages := []models.Package{}
	for rows.Next() {
		var p models.Package
		var author models.User

		err := rows.Scan(
			&p.ID, &p.Slug, &p.DisplayName, &p.Description, &p.Type, &p.Status,
			&p.RepositoryURL, &p.License, &p.AuthorID, &p.CreatedAt, &p.UpdatedAt,
			&p.ViewCount, &p.BookmarkCount,
			&author.Username, &author.Alias, &author.AvatarURL,
		)
		if err != nil {
			fmt.Println("Iter row err", err)
			continue
		}
		p.Author = &author
		packages = append(packages, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packages)
}
