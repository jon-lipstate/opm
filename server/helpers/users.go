package helpers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"opm/db"
	"opm/models"
)

// FindOrCreateUser finds an existing user or creates a new one
func FindOrCreateUser(ctx context.Context, provider string, providerID string, username string, displayName string, avatarURL string) (*models.User, error) {
	// Create a new context with a longer timeout for database operations
	dbCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	var user models.User
	var query string
	var args []interface{}

	// Check which provider we're using
	switch provider {
	case "github":
		query = `SELECT id, github_id, discord_id, username, alias, display_name, avatar_url, created_at, updated_at 
				 FROM users WHERE github_id = $1`
		args = []interface{}{providerID}
	case "discord":
		query = `SELECT id, github_id, discord_id, username, alias, display_name, avatar_url, created_at, updated_at 
				 FROM users WHERE discord_id = $1`
		args = []interface{}{providerID}
	default:
		return nil, fmt.Errorf("invalid provider: %s", provider)
	}

	// Try to find existing user
	err := db.QueryRow(dbCtx, query, args...).Scan(
		&user.ID,
		&user.GitHubID,
		&user.DiscordID,
		&user.Username,
		&user.Alias,
		&user.DisplayName,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == nil {
		// User exists, update their info
		updateQuery := `UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2`
		_, err = db.Exec(dbCtx, updateQuery, avatarURL, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
		user.AvatarURL = &avatarURL
		return &user, nil
	}

	if err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// User doesn't exist, create new one
	// First, ensure username is unique
	baseUsername := strings.ToLower(username)
	uniqueUsername := baseUsername
	counter := 1

	for {
		var exists bool
		err = db.QueryRow(dbCtx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", uniqueUsername).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
		if !exists {
			break
		}
		uniqueUsername = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++
	}

	// Generate a unique alias
	baseAlias := GenerateAlias(uniqueUsername)
	uniqueAlias := baseAlias
	counter = 1

	for {
		var exists bool
		err = db.QueryRow(dbCtx, "SELECT EXISTS(SELECT 1 FROM users WHERE alias = $1)", uniqueAlias).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check alias: %w", err)
		}
		if !exists {
			break
		}
		uniqueAlias = fmt.Sprintf("%s-%d", baseAlias, counter)
		counter++
	}

	// Create new user
	var insertQuery string
	if provider == "github" {
		insertQuery = `INSERT INTO users (github_id, username, alias, display_name, avatar_url) 
					   VALUES ($1, $2, $3, $4, $5) 
					   RETURNING id, github_id, discord_id, username, alias, display_name, avatar_url, created_at, updated_at`
		args = []interface{}{providerID, uniqueUsername, uniqueAlias, displayName, avatarURL}
	} else {
		insertQuery = `INSERT INTO users (discord_id, username, alias, display_name, avatar_url) 
					   VALUES ($1, $2, $3, $4, $5) 
					   RETURNING id, github_id, discord_id, username, alias, display_name, avatar_url, created_at, updated_at`
		args = []interface{}{providerID, uniqueUsername, uniqueAlias, displayName, avatarURL}
	}

	err = db.QueryRow(dbCtx, insertQuery, args...).Scan(
		&user.ID,
		&user.GitHubID,
		&user.DiscordID,
		&user.Username,
		&user.Alias,
		&user.DisplayName,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}