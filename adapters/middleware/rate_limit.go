package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateBucket struct {
	count   int
	resetAt time.Time
}

// RateLimit constrains anonymous, high-risk endpoints such as login and
// registration. It deliberately uses the direct remote address because proxy
// forwarding headers must only be trusted after a deployment configures its
// trusted proxy list.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	var (
		mu      sync.Mutex
		buckets = make(map[string]rateBucket)
	)
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
		if err != nil || ip == "" {
			ip = c.ClientIP()
		}
		now := time.Now()
		mu.Lock()
		bucket := buckets[ip]
		if bucket.resetAt.IsZero() || !now.Before(bucket.resetAt) {
			bucket = rateBucket{resetAt: now.Add(window)}
		}
		bucket.count++
		buckets[ip] = bucket
		allowed := bucket.count <= limit
		retryAfter := int(time.Until(bucket.resetAt).Seconds()) + 1
		mu.Unlock()
		if !allowed {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "demasiados intentos; intenta de nuevo más tarde"})
			c.Abort()
			return
		}
		c.Next()
	}
}
