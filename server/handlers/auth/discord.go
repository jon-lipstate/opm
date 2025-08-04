package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"opm/config"
	"opm/helpers"
)

var discordEndpoint = oauth2.Endpoint{
	AuthURL:  "https://discord.com/api/oauth2/authorize",
	TokenURL: "https://discord.com/api/oauth2/token",
}

// DiscordLogin initiates the Discord OAuth flow
func DiscordLogin(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Construct full redirect URL
		redirectURL := cfg.Host
		// In production, don't add port if HOST already includes the full URL
		if cfg.Env == "development" && cfg.Port != "" && cfg.Port != "80" && cfg.Port != "443" {
			redirectURL = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		}
		redirectURL = redirectURL + "/" + cfg.DiscordRedirectURL

		oauthConfig := &oauth2.Config{
			ClientID:     cfg.DiscordClientID,
			ClientSecret: cfg.DiscordClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"identify", "email"},
			Endpoint:     discordEndpoint,
		}

		// Generate state token using JWT secret
		state := helpers.GenerateState(cfg.JWTSecret)
		
		url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// DiscordCallback handles the Discord OAuth callback
func DiscordCallback(cfg *config.Config) http.HandlerFunc {
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
		redirectURL = redirectURL + "/" + cfg.DiscordRedirectURL

		oauthConfig := &oauth2.Config{
			ClientID:     cfg.DiscordClientID,
			ClientSecret: cfg.DiscordClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"identify", "email"},
			Endpoint:     discordEndpoint,
		}

		// Exchange code for token
		token, err := oauthConfig.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		// Get user info from Discord
		client := oauthConfig.Client(r.Context(), token)
		resp, err := client.Get("https://discord.com/api/users/@me")
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var discordUser struct {
			ID            string `json:"id"`
			Username      string `json:"username"`
			Discriminator string `json:"discriminator"`
			Avatar        string `json:"avatar"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&discordUser); err != nil {
			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
			return
		}

		// TODO: Create or update user in database
		// TODO: Generate JWT token
		// TODO: Set cookie and redirect to frontend

		// For now, just return success
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Discord login successful",
			"user": map[string]interface{}{
				"discord_id": discordUser.ID,
				"username":   discordUser.Username,
				"avatar":     discordUser.Avatar,
			},
		})
	}
}