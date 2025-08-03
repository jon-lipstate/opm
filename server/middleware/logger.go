package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger middleware logs HTTP requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap response writer
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		
		// Process request
		next.ServeHTTP(rw, r)
		
		// Log request
		duration := time.Since(start)
		log.Printf(
			"%s %s %s %d %v",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			rw.status,
			duration,
		)
	})
}