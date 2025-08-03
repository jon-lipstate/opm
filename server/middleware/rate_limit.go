package middleware

import (
	"net/http"
	"sync"
	"time"
)

// Simple in-memory rate limiter
type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

// RateLimit creates a rate limiting middleware
func RateLimit(limit string, window string) func(http.Handler) http.Handler {
	// Parse limit and window
	// TODO: Proper parsing
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    100, // Default to 100 requests
		window:   time.Minute, // Default to 1 minute
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			
			rl.mu.Lock()
			defer rl.mu.Unlock()
			
			now := time.Now()
			cutoff := now.Add(-rl.window)
			
			// Get existing requests for this IP
			requests := rl.requests[ip]
			
			// Filter out old requests
			var validRequests []time.Time
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}
			
			// Check if limit exceeded
			if len(validRequests) >= rl.limit {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			
			// Add current request
			validRequests = append(validRequests, now)
			rl.requests[ip] = validRequests
			
			next.ServeHTTP(w, r)
		})
	}
}