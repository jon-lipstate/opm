package packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type GitHubRepo struct {
	Name        string  `json:"name"`
	FullName    string  `json:"full_name"`
	Description string  `json:"description"`
	License     *struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	} `json:"license"`
	Homepage string `json:"homepage"`
	Topics   []string `json:"topics"`
}

// GetRepositoryMetadata fetches metadata from a repository URL
func GetRepositoryMetadata(w http.ResponseWriter, r *http.Request) {
	repoURL := r.URL.Query().Get("url")
	if repoURL == "" {
		http.Error(w, "Repository URL is required", http.StatusBadRequest)
		return
	}

	// Parse the URL
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		http.Error(w, "Invalid repository URL", http.StatusBadRequest)
		return
	}

	// Extract owner and repo from GitHub URL
	if parsedURL.Host != "github.com" && parsedURL.Host != "www.github.com" {
		http.Error(w, "Only GitHub repositories are supported for metadata extraction", http.StatusBadRequest)
		return
	}

	// Remove leading/trailing slashes and .git extension
	path := strings.Trim(parsedURL.Path, "/")
	path = strings.TrimSuffix(path, ".git")
	
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid GitHub repository URL format", http.StatusBadRequest)
		return
	}

	owner := parts[0]
	repo := parts[1]

	// Fetch from GitHub API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Add headers for GitHub API
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "OPM-Package-Registry")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch repository data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			http.Error(w, "Repository not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch repository data", http.StatusInternalServerError)
		}
		return
	}

	var githubRepo GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&githubRepo); err != nil {
		http.Error(w, "Failed to parse repository data", http.StatusInternalServerError)
		return
	}

	// Prepare response with suggested values
	response := map[string]interface{}{
		"display_name": githubRepo.Name,
		"description":  githubRepo.Description,
		"topics":       githubRepo.Topics,
	}

	// Add license if available
	if githubRepo.License != nil {
		response["license"] = githubRepo.License.Name
	}

	// Generate a slug from the repo name
	slug := strings.ToLower(githubRepo.Name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove any characters that aren't alphanumeric or hyphens
	cleanSlug := ""
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			cleanSlug += string(r)
		}
	}
	response["slug"] = cleanSlug

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}