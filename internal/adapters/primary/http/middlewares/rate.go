package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	burst    int
}

func newIPLimiter(r rate.Limit, burst int) *ipLimiter {
	return &ipLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		burst:    burst,
	}
}

func (l *ipLimiter) allow(ip string) bool {
	l.mu.Lock()
	lim, exists := l.limiters[ip]
	if !exists {
		lim = rate.NewLimiter(l.r, l.burst)
		l.limiters[ip] = lim
	}
	l.mu.Unlock()
	return lim.Allow()
}

func AuthRateLimiter() gin.HandlerFunc {
	limiter := newIPLimiter(2, 10)

	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many attempts, please slow down",
			})
			return
		}
		c.Next()
	}
}

func APIRateLimiter() gin.HandlerFunc {
	limiter := newIPLimiter(50, 100)

	return func(c *gin.Context) {
		if !limiter.allow(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
