package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"opm/config"
	"opm/helpers"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GitHubLogin initiates the GitHub OAuth flow
func GitHubLogin(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Construct full redirect URL
		redirectURL := cfg.Host
		// In production, don't add port if HOST already includes the full URL
		if cfg.Env == "development" && cfg.Port != "" && cfg.Port != "80" && cfg.Port != "443" {
			redirectURL = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		}
		redirectURL = redirectURL + "/" + cfg.GitHubRedirectURL

		oauthConfig := &oauth2.Config{
			ClientID:     cfg.GitHubClientID,
			ClientSecret: cfg.GitHubClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		}

		// Generate state token using JWT secret
		state := helpers.GenerateState(cfg.JWTSecret)

		url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// GitHubCallback handles the GitHub OAuth callback
func GitHubCallback(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify state parameter
		state := r.URL.Query().Get("state")
		if !helpers.ValidateState(state, cfg.JWTSecret) {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code parameter", http.StatusBadRequest)
			return
		}

		// Construct full redirect URL
		redirectURL := cfg.Host
		// In production, don't add port if HOST already includes the full URL
		if cfg.Env == "development" && cfg.Port != "" && cfg.Port != "80" && cfg.Port != "443" {
			redirectURL = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		}
		redirectURL = redirectURL + "/" + cfg.GitHubRedirectURL

		oauthConfig := &oauth2.Config{
			ClientID:     cfg.GitHubClientID,
			ClientSecret: cfg.GitHubClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		}

		// Exchange code for token
		token, err := oauthConfig.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		// Get user info from GitHub
		client := oauthConfig.Client(r.Context(), token)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var githubUser struct {
			ID        int64  `json:"id"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
			return
		}

		// Create or update user in database
		displayName := githubUser.Name
		if displayName == "" {
			displayName = githubUser.Login
		}

		user, err := helpers.FindOrCreateUser(
			r.Context(),
			"github",
			fmt.Sprintf("%d", githubUser.ID),
			githubUser.Login,
			displayName,
			githubUser.AvatarURL,
		)
		if err != nil {
			// Log the actual error for debugging
			fmt.Printf("Failed to create/update user: %v\n", err)
			http.Error(w, "Failed to create/update user", http.StatusInternalServerError)
			return
		}

		// Generate JWT token
		tokenString, err := helpers.GenerateJWT(user.ID, cfg.JWTSecret)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		var domain string
		var secure bool
		switch os.Getenv("ENV") {
		case "development":
			domain = "localhost"
			secure = false
		case "production":
			domain = ".pkg-odin.org"  // Allow cookie across subdomains
			secure = true
		default:
			fmt.Println("INVALID ENVIRONMENT")
			domain = "localhost"
			secure = false
		}

		// Set cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Path:     "/",
			Domain:   domain,
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   7 * 24 * 60 * 60, // 7 days
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		// Redirect to frontend
		redirectURL = fmt.Sprintf("%s?auth=success", cfg.FrontendURL)
		fmt.Println("REDIRECT:", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}
