package packages

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"opm/db"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// GetPackageReadme fetches the README content from the package's repository
func GetPackageReadme(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get package ID from query parameter
	packageIDStr := r.URL.Query().Get("package_id")
	if packageIDStr == "" {
		http.Error(w, "Missing package_id parameter", http.StatusBadRequest)
		return
	}

	var packageID int
	if _, err := fmt.Sscanf(packageIDStr, "%d", &packageID); err != nil {
		http.Error(w, "Invalid package_id", http.StatusBadRequest)
		return
	}

	// Get package repository URL
	var repositoryURL string
	err := db.Conn.QueryRow(ctx,
		"SELECT repository_url FROM packages WHERE id = $1",
		packageID,
	).Scan(&repositoryURL)
	if err == pgx.ErrNoRows {
		http.Error(w, "Package not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to find package", http.StatusInternalServerError)
		return
	}

	// Parse repository URL and fetch README
	readmeContent, err := fetchReadmeFromRepo(repositoryURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch README: %v", err), http.StatusInternalServerError)
		return
	}

	// Return README content
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"content": readmeContent,
		"format":  "markdown",
	})
}

// fetchReadmeFromRepo fetches README content from various repository providers
func fetchReadmeFromRepo(repoURL string) (string, error) {
	// Clean up the URL
	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.TrimSuffix(repoURL, "/")

	// GitHub repositories
	if strings.Contains(repoURL, "github.com") {
		return fetchGitHubReadme(repoURL)
	}

	// GitLab repositories
	if strings.Contains(repoURL, "gitlab.com") {
		return fetchGitLabReadme(repoURL)
	}

	// Codeberg repositories
	if strings.Contains(repoURL, "codeberg.org") {
		return fetchCodebergReadme(repoURL)
	}

	return "", fmt.Errorf("unsupported repository provider")
}

// fetchGitHubReadme fetches README from GitHub
func fetchGitHubReadme(repoURL string) (string, error) {
	// Extract owner and repo from URL
	// Example: https://github.com/owner/repo -> owner/repo
	parts := strings.Split(repoURL, "github.com/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid GitHub URL")
	}

	repoPath := strings.TrimPrefix(parts[1], "/")

	// Try different README filenames
	readmeFiles := []string{"README.md", "readme.md", "README.MD", "Readme.md", "README", "readme"}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, filename := range readmeFiles {
		// Use GitHub raw content URL
		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/main/%s", repoPath, filename)

		// Try main branch first
		content, err := fetchFromURL(client, rawURL)
		if err == nil {
			return content, nil
		}

		// Try master branch
		rawURL = fmt.Sprintf("https://raw.githubusercontent.com/%s/master/%s", repoPath, filename)
		content, err = fetchFromURL(client, rawURL)
		if err == nil {
			return content, nil
		}
	}

	return "", fmt.Errorf("README not found")
}

// fetchGitLabReadme fetches README from GitLab
func fetchGitLabReadme(repoURL string) (string, error) {
	// Extract project path from URL
	parts := strings.Split(repoURL, "gitlab.com/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid GitLab URL")
	}

	projectPath := strings.TrimPrefix(parts[1], "/")

	// Try different README filenames
	readmeFiles := []string{"README.md", "readme.md", "README.MD", "Readme.md", "README", "readme"}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, filename := range readmeFiles {
		// Use GitLab raw content URL
		rawURL := fmt.Sprintf("https://gitlab.com/%s/-/raw/main/%s", projectPath, filename)

		// Try main branch first
		content, err := fetchFromURL(client, rawURL)
		if err == nil {
			return content, nil
		}

		// Try master branch
		rawURL = fmt.Sprintf("https://gitlab.com/%s/-/raw/master/%s", projectPath, filename)
		content, err = fetchFromURL(client, rawURL)
		if err == nil {
			return content, nil
		}
	}

	return "", fmt.Errorf("README not found")
}

// fetchCodebergReadme fetches README from Codeberg
func fetchCodebergReadme(repoURL string) (string, error) {
	// Extract owner and repo from URL
	parts := strings.Split(repoURL, "codeberg.org/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid Codeberg URL")
	}

	repoPath := strings.TrimPrefix(parts[1], "/")

	// Try different README filenames
	readmeFiles := []string{"README.md", "readme.md", "README.MD", "Readme.md", "README", "readme"}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, filename := range readmeFiles {
		// Use Codeberg raw content URL
		rawURL := fmt.Sprintf("https://codeberg.org/%s/raw/branch/main/%s", repoPath, filename)

		// Try main branch first
		content, err := fetchFromURL(client, rawURL)
		if err == nil {
			return content, nil
		}

		// Try master branch
		rawURL = fmt.Sprintf("https://codeberg.org/%s/raw/branch/master/%s", repoPath, filename)
		content, err = fetchFromURL(client, rawURL)
		if err == nil {
			return content, nil
		}
	}

	return "", fmt.Errorf("README not found")
}

// fetchFromURL fetches content from a URL
func fetchFromURL(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	// Limit README size to 1MB
	limitedReader := io.LimitReader(resp.Body, 1024*1024)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
