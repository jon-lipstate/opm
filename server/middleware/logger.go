package middleware

import (
	"context"
	"fmt"
	"net/http"
	"opm/logger"
	"strings"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code and user info
type responseWriter struct {
	http.ResponseWriter
	status int
	userID *int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// SetUserID allows auth middleware to set the user ID for logging
func (rw *responseWriter) SetUserID(userID int) {
	rw.userID = &userID
}

// Logger middleware logs HTTP requests
func Logger(next http.Handler) http.Handler {
	mainLogger := logger.MainLogger
	securityLogger := logger.SecurityLogger
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Get client IP
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.Header.Get("X-Real-IP")
		}
		if clientIP == "" {
			clientIP = r.RemoteAddr
		} else {
			// Take first IP if multiple
			clientIP = strings.Split(clientIP, ",")[0]
		}
		
		// Skip noisy endpoints
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		
		// Wrap response writer
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		
		// Pass wrapped writer to context for auth middleware
		ctx := context.WithValue(r.Context(), "responseWriter", rw)
		
		// Process request
		next.ServeHTTP(rw, r.WithContext(ctx))
		
		// Build user identifier
		var userIdentifier string
		if rw.userID != nil {
			userIdentifier = fmt.Sprintf("%s|u:%d", clientIP, *rw.userID)
		} else {
			userIdentifier = clientIP
		}
		
		// Log based on status code
		duration := time.Since(start)
		logLine := fmt.Sprintf("[%s] %s %d %s in %v", userIdentifier, r.Method, rw.status, r.URL.Path, duration)
		
		if rw.status >= 400 && rw.status < 500 {
			// Client errors go to security log
			securityLogger.Println(logLine)
			if rw.status == 404 {
				securityLogger.Printf("[%s] 404 Not Found: %s %s", clientIP, r.Method, r.URL.Path)
			}
		} else {
			// Normal requests go to main log
			mainLogger.Println(logLine)
		}
	})
}