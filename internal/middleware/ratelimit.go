package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple in-memory rate limiting mechanism
type RateLimiter struct {
	// Map to store request counts per IP
	requests map[string][]time.Time
	// Duration for the rate limiting window
	window time.Duration
	// Maximum number of requests allowed in the window
	limit int
	// Mutex for concurrent access
	mu sync.Mutex
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		limit:    limit,
	}
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Initialize if this is the first request from this IP
	if _, exists := rl.requests[ip]; !exists {
		rl.requests[ip] = []time.Time{}
	}

	// Clean up old requests outside the current window
	cutoff := now.Add(-rl.window)
	var validRequests []time.Time

	for _, reqTime := range rl.requests[ip] {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	rl.requests[ip] = validRequests

	// Check if we've hit the limit
	if len(rl.requests[ip]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

// GetClientIP extracts the client IP from the request
func GetClientIP(r *http.Request) string {
	// First try the X-Forwarded-For header
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, use the first one
		if ips := net.ParseIP(ip); ips != nil {
			return ips.String()
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// RateLimit creates a middleware that applies rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := GetClientIP(r)

		if !rl.Allow(clientIP) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			http.Error(w, `{"error":"Too many requests, please try again later"}`, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimitFunc creates a middleware function that applies rate limiting
func (rl *RateLimiter) MiddlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := GetClientIP(r)

		if !rl.Allow(clientIP) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			http.Error(w, `{"error":"Too many requests, please try again later"}`, http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
