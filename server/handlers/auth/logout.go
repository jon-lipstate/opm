package auth

import (
	"net/http"
	"os"
	"time"
)

// Logout handles user logout
func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var domain string
		var secure bool
		switch os.Getenv("ENV") {
		case "development":
			domain = "localhost"
			secure = false
		case "production":
			domain = ".pkg-odin.org"  // Match login cookie domain
			secure = true
		default:
			domain = "localhost"
			secure = false
		}

		// Clear the auth cookie - try multiple domains to handle old cookies
		domains := []string{domain}
		if os.Getenv("ENV") == "production" {
			// Also clear any old cookies that might exist
			domains = append(domains, "api.pkg-odin.org", "pkg-odin.org", ".api.pkg-odin.org")
		}

		for _, d := range domains {
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    "",
				Path:     "/",
				Domain:   d,
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   secure,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(-1 * time.Hour),
			})
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Logged out successfully"}`))
	}
}