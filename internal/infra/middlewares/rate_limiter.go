package middlewares

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiterConfig defines limit and period
type RateLimiterConfig struct {
	Limit  int
	Period time.Duration
}

// visitor holds the limiter for each client IP
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

// cleanup visitors periodically to prevent memory leak
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// getVisitor returns the rate limiter for a given IP
func getVisitor(ip string, cfg RateLimiterConfig) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Every(cfg.Period/time.Duration(cfg.Limit)), cfg.Limit)
		visitors[ip] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// RateLimiter is a middleware that provides rate limiting functionality.
func RateLimiter(next http.Handler, cfg RateLimiterConfig) http.Handler {
	// Start cleanup goroutine
	go cleanupVisitors()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // This will not be correct if behind a proxy
		limiter := getVisitor(ip, cfg)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded. Try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}