package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

// NewRateLimiterMiddleware creates a Gin middleware with rate limiting
func NewRateLimiterMiddleware(cfg RateLimiterConfig) gin.HandlerFunc {
	// Start cleanup goroutine
	go cleanupVisitors()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitor(ip, cfg)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			return
		}

		c.Next()
	}
}
