package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID               int        `json:"id"`
	GitHubID         *string    `json:"github_id,omitempty"`
	DiscordID        *string    `json:"discord_id,omitempty"`
	Username         string     `json:"username"`
	Slug             string     `json:"slug"`
	DisplayName      *string    `json:"display_name,omitempty"`
	AvatarURL        *string    `json:"avatar_url,omitempty"`
	IsModerator      bool       `json:"is_moderator"`
	Reputation       int        `json:"reputation"`
	ReputationRank   string     `json:"reputation_rank"`
	DiscordVerified  bool       `json:"discord_verified"`
	GitHubVerified   bool       `json:"github_verified"`
	IsBanned         bool       `json:"is_banned"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Package represents a package in the registry
type Package struct {
	ID            int           `json:"id"`
	Slug          string        `json:"slug"`  // URL-safe name
	DisplayName   string        `json:"display_name"`
	Description   string        `json:"description"`
	Type          PackageType   `json:"type"`
	Status        PackageStatus `json:"status"`
	RepositoryURL string        `json:"repository_url"`
	License       *string       `json:"license,omitempty"`
	AuthorID      int           `json:"author_id"`
	Author        *User         `json:"author,omitempty"`
	Tags          []Tag         `json:"tags,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	
	// Stats
	ViewCount     int  `json:"view_count"`
	BookmarkCount int  `json:"bookmark_count"`
	
	// Computed fields
	IsBookmarked       bool `json:"is_bookmarked"`
	ActiveReportsCount int  `json:"active_reports_count"`
}

// Tag represents a tag that can be applied to packages
type Tag struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	AddedBy    *int      `json:"added_by,omitempty"`
	UsageCount int       `json:"usage_count"`
	CreatedAt  time.Time `json:"created_at"`
	
	// For package-tag relation
	NetScore  int `json:"net_score"`  // Sum of all votes
	UserVote  int `json:"user_vote"`  // Current user's vote (-1, 0, 1)
}

// Bookmark represents a user's bookmarked package
type Bookmark struct {
	UserID    int       `json:"user_id"`
	PackageID int       `json:"package_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePackageInput represents the input for creating a new package
type CreatePackageInput struct {
	Slug          string        `json:"slug" validate:"required,min=2,max=100,slug"`
	DisplayName   string        `json:"display_name" validate:"required,min=2,max=255"`
	Description   string        `json:"description" validate:"required,min=10,max=1000"`
	Type          PackageType   `json:"type" validate:"required,oneof=library project"`
	Status        PackageStatus `json:"status" validate:"required,oneof=in_work ready archived abandoned"`
	RepositoryURL string        `json:"repository_url" validate:"required,url"`
	License       *string       `json:"license,omitempty" validate:"omitempty,max=100"`
	TagIDs        []int         `json:"tag_ids"`
}

// UpdatePackageInput represents the input for updating a package
type UpdatePackageInput struct {
	DisplayName   *string        `json:"display_name,omitempty" validate:"omitempty,min=2,max=255"`
	Description   *string        `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Type          *PackageType   `json:"type,omitempty" validate:"omitempty,oneof=library project"`
	Status        *PackageStatus `json:"status,omitempty" validate:"omitempty,oneof=in_work ready"`
	RepositoryURL *string        `json:"repository_url,omitempty" validate:"omitempty,url"`
	License       *string        `json:"license,omitempty" validate:"omitempty,max=100"`
	TagIDs        []int          `json:"tag_ids,omitempty"`
}

// PackageFilter represents filters for querying packages
type PackageFilter struct {
	Type     *PackageType   `json:"type,omitempty"`
	Status   *PackageStatus `json:"status,omitempty"`
	Tags     []string       `json:"tags,omitempty"`  // Tag names, not slugs
	AuthorID *int           `json:"author_id,omitempty"`
	Search   *string        `json:"search,omitempty"`
	Limit    int            `json:"limit,omitempty"`
	Offset   int            `json:"offset,omitempty"`
}

// AuthUser represents the authenticated user stored in context
type AuthUser struct {
	UserID int
	Token  string
}

// PackageView represents a view of a package
type PackageView struct {
	PackageID int       `json:"package_id"`
	UserID    *int      `json:"user_id,omitempty"`
	ViewedAt  time.Time `json:"viewed_at"`
}

// UpdateUserInput represents the input for updating a user profile
type UpdateUserInput struct {
	Slug        *string `json:"slug,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

// Flag represents a moderation flag on a package
type Flag struct {
	ID         int        `json:"id"`
	PackageID  int        `json:"package_id"`
	UserID     int        `json:"user_id"`
	Reason     string     `json:"reason"`
	Details    *string    `json:"details,omitempty"`
	Status     string     `json:"status"`  // pending, reviewed, resolved, dismissed
	ResolvedBy *int       `json:"resolved_by,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TagVote represents a user's vote on a package tag
type TagVote struct {
	ID        int       `json:"id"`
	PackageID int       `json:"package_id"`
	TagID     int       `json:"tag_id"`
	UserID    int       `json:"user_id"`
	VoteValue int       `json:"vote_value"` // -10 to +10
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateFlagInput represents the input for flagging a package
type CreateFlagInput struct {
	Reason  string  `json:"reason" validate:"required"`
	Details *string `json:"details,omitempty"`
}

// TagVoteInput represents the input for voting on a tag
type TagVoteInput struct {
	Vote int `json:"vote" validate:"required,min=-1,max=1"` // -1, 0, or 1
}