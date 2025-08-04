package packages

import (
	"encoding/json"
	"net/http"
	"opm/db"
	"opm/helpers"
	"opm/logger"
	"opm/middleware"
	"opm/models"
)

// Search performs full-text search on packages
func Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get search query
	searchQuery, ok := helpers.RequiredParamString(r, w, "q")
	if !ok {
		return
	}

	// Pagination
	limit := 20
	offset := 0
	if l, hasLimit := helpers.OptionalParamInt(r, "limit"); hasLimit {
		limit = *l
	}

	if o, hasOffset := helpers.OptionalParamInt(r, "offset"); hasOffset {
		offset = *o
	}

	// Build search query
	query := `
		SELECT p.id, p.slug, p.display_name, p.description, p.type, p.status,
		       p.repository_url, p.author_id, p.created_at, p.updated_at,
		       (SELECT COUNT(*) FROM package_views WHERE package_id = p.id) as view_count,
		       (SELECT COUNT(*) FROM bookmarks WHERE package_id = p.id) as bookmark_count,
		       u.username, u.slug, u.display_name, u.avatar_url,
		       (SELECT COUNT(*) FROM flags WHERE package_id = p.id AND status = 'pending') as active_reports_count,
		       ts_rank(p.search_vector, plainto_tsquery('english', $1)) as rank
		FROM packages p
		JOIN users u ON p.author_id = u.id
		WHERE p.search_vector @@ plainto_tsquery('english', $1)
		ORDER BY rank DESC, p.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := db.Conn.Query(ctx, query, searchQuery, limit, offset)
	if err != nil {
		logger.MainLogger.Printf("Search query error - Query: %s, Error: %v", searchQuery, err)
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	packages := []models.Package{}
	for rows.Next() {
		var p models.Package
		var author models.User
		var activeReportsCount int
		var rank float32

		err := rows.Scan(
			&p.ID, &p.Slug, &p.DisplayName, &p.Description, &p.Type, &p.Status,
			&p.RepositoryURL, &p.AuthorID, &p.CreatedAt, &p.UpdatedAt,
			&p.ViewCount, &p.BookmarkCount,
			&author.Username, &author.Slug, &author.DisplayName, &author.AvatarURL,
			&activeReportsCount,
			&rank,
		)
		p.ActiveReportsCount = activeReportsCount
		if err != nil {
			logger.MainLogger.Printf("Failed to scan search result: %v", err)
			continue
		}
		p.Author = &author
		packages = append(packages, p)
	}

	// Get tags for each package
	for i := range packages {
		tags, err := getPackageTags(ctx, packages[i].ID, 0)
		if err == nil {
			packages[i].Tags = tags
		} else {
			logger.MainLogger.Printf("Failed to get tags for package %d: %v", packages[i].ID, err)
		}
	}

	// Check bookmarks if user is authenticated
	if authUser, ok := middleware.GetAuthUser(ctx); ok {
		for i := range packages {
			packages[i].IsBookmarked = checkBookmark(ctx, authUser.UserID, packages[i].ID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packages)
}
