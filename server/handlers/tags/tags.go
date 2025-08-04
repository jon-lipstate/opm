package tags

import (
	"encoding/json"
	"net/http"
	"opm/db"
	"opm/logger"
	"opm/models"
	"strconv"
)

// List returns all tags or search for tags
func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get query parameters
	search := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")

	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	query := `
			SELECT id, name, usage_count, created_at
			FROM tags
			WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += " AND search_vector @@ plainto_tsquery('english', $" + strconv.Itoa(argIndex) + ")"
		args = append(args, search)
		argIndex++
	}

	query += " ORDER BY usage_count DESC, name"
	query += " LIMIT $" + strconv.Itoa(argIndex)
	args = append(args, limit)

	rows, err := db.Conn.Query(ctx, query, args...)
	if err != nil {
		logger.MainLogger.Printf("Failed to fetch tags - Query: %s, Args: %v, Error: %v", query, args, err)
		http.Error(w, "Failed to fetch tags", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		var t models.Tag
		err := rows.Scan(&t.ID, &t.Name, &t.UsageCount, &t.CreatedAt)
		if err != nil {
			logger.MainLogger.Printf("Failed to scan tag: %v", err)
			continue
		}
		tags = append(tags, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}
