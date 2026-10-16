package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"satset-app/internal/handler"
)

type rateLimitEntry struct {
	attempts    int
	lockedUntil time.Time
	firstSeen   time.Time
}

// MemoryRateLimiter implements an in-memory rate limiter.
type MemoryRateLimiter struct {
	mu              sync.Mutex
	entries         map[string]*rateLimitEntry
	maxAttempts     int
	lockoutDuration time.Duration
	windowDuration  time.Duration
}

// NewMemoryRateLimiter creates a new rate limiter.
func NewMemoryRateLimiter(maxAttempts int, lockoutDuration, windowDuration time.Duration) *MemoryRateLimiter {
	limiter := &MemoryRateLimiter{
		entries:         make(map[string]*rateLimitEntry),
		maxAttempts:     maxAttempts,
		lockoutDuration: lockoutDuration,
		windowDuration:  windowDuration,
	}

	// Background routine to cleanup old entries every 5 minutes
	go limiter.cleanupLoop()

	return limiter
}

func (l *MemoryRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for k, v := range l.entries {
			if now.After(v.lockedUntil) && now.Sub(v.firstSeen) > l.windowDuration {
				delete(l.entries, k)
			}
		}
		l.mu.Unlock()
	}
}

// IsAllowed checks if the key is allowed and increments attempt count.
func (l *MemoryRateLimiter) IsAllowed(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	entry, exists := l.entries[key]
	if !exists {
		l.entries[key] = &rateLimitEntry{
			attempts:  1,
			firstSeen: now,
		}
		return true, 0
	}

	if now.Before(entry.lockedUntil) {
		return false, entry.lockedUntil.Sub(now)
	}

	// Reset if outside sliding window
	if now.Sub(entry.firstSeen) > l.windowDuration {
		entry.attempts = 1
		entry.firstSeen = now
		entry.lockedUntil = time.Time{}
		return true, 0
	}

	entry.attempts++
	if entry.attempts > l.maxAttempts {
		entry.lockedUntil = now.Add(l.lockoutDuration)
		return false, l.lockoutDuration
	}

	return true, 0
}

// Reset clears the rate limit for a key (e.g., upon successful login).
func (l *MemoryRateLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key)
}

// LoginRateLimitMiddleware applies rate limiting per client IP to login requests.
func LoginRateLimitMiddleware(limiter *MemoryRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		allowed, waitTime := limiter.IsAllowed(ip)
		if !allowed {
			mins := int(waitTime.Minutes()) + 1
			handler.RespondTooManyRequests(c, fmt.Sprintf("Terlalu banyak percobaan gagal. Silakan coba lagi dalam %d menit.", mins))
			c.Abort()
			return
		}
		c.Next()
		if c.Writer.Status() == http.StatusOK {
			limiter.Reset(ip)
		}
	}
}
