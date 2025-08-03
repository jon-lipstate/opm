package packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opm/db"
	"opm/middleware"
	"opm/models"
	"strconv"
)

// Search performs full-text search on packages
func Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get search query
	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		http.Error(w, "Search query required", http.StatusBadRequest)
		return
	}

	// Pagination
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	fmt.Println(searchQuery, limit, offset)

	// Build search query
	query := `
		SELECT p.id, p.slug, p.display_name, p.description, p.type, p.status,
		       p.repository_url, p.author_id, p.created_at, p.updated_at,
		       (SELECT COUNT(*) FROM package_views WHERE package_id = p.id) as view_count,
		       (SELECT COUNT(*) FROM bookmarks WHERE package_id = p.id) as bookmark_count,
		       u.username, u.alias, u.display_name, u.avatar_url,
		       (SELECT COUNT(*) FROM flags WHERE package_id = p.id AND status = 'pending') as active_reports_count,
		       ts_rank(p.search_vector, plainto_tsquery('english', $1)) as rank
		FROM packages p
		JOIN users u ON p.author_id = u.id
		WHERE p.search_vector @@ plainto_tsquery('english', $1)
		ORDER BY rank DESC, p.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := db.Conn.Query(ctx, query, searchQuery, limit, offset)
	if err != nil {
		fmt.Printf("Search query error: %v\n", err)
		fmt.Printf("Search term: %s\n", searchQuery)
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
			&author.Username, &author.Alias, &author.DisplayName, &author.AvatarURL,
			&activeReportsCount,
			&rank,
		)
		p.ActiveReportsCount = activeReportsCount
		if err != nil {
			fmt.Printf("Search scan error: %v\n", err)
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
			fmt.Println("TAG ERR", err)
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
