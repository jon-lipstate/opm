package auth

import (
	"net/http"
	"time"
)

// Logout handles user logout
func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear the auth cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(-1 * time.Hour),
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Logged out successfully"}`))
	}
}