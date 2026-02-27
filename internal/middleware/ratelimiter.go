package middleware

import (
	"net/http"
	"sync"
	"time"
)

// visitor tracks requests for a single key inside the current time window.
type visitor struct {
	windowStart time.Time
	count       int
}

// RateLimiter implements a simple fixed-window limiter in memory.
type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	storage map[string]visitor
}

// NewRateLimiter creates a limiter with request limit per time window.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:   limit,
		window:  window,
		storage: make(map[string]visitor),
	}
}

// Allow returns true when the key can perform one more request in window.
func (rl *RateLimiter) Allow(key string, now time.Time) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.storage[key]
	if !exists || now.Sub(entry.windowStart) >= rl.window {
		rl.storage[key] = visitor{windowStart: now, count: 1}
		return true
	}

	if entry.count >= rl.limit {
		return false
	}

	entry.count++
	rl.storage[key] = entry
	return true
}

// Handler wraps an endpoint and rejects requests once the limit is hit.
func (rl *RateLimiter) Handler(next http.HandlerFunc, keyFn func(*http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow(keyFn(r), time.Now()) {
			http.Error(w, "Too many requests. Please retry shortly", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
