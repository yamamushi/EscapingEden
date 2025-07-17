package edenutil

import (
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mutex      sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if an action is allowed under the rate limit
func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Refill tokens based on elapsed time
	if elapsed >= rl.refillRate {
		tokensToAdd := int(elapsed / rl.refillRate)
		rl.tokens = min(rl.maxTokens, rl.tokens+tokensToAdd)
		rl.lastRefill = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// ConnectionRateLimiter manages rate limiting per IP address
type ConnectionRateLimiter struct {
	limiters   map[string]*RateLimiter
	mutex      sync.RWMutex
	maxTokens  int
	refillRate time.Duration
}

// NewConnectionRateLimiter creates a new connection rate limiter
func NewConnectionRateLimiter(maxTokens int, refillRate time.Duration) *ConnectionRateLimiter {
	return &ConnectionRateLimiter{
		limiters:   make(map[string]*RateLimiter),
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}
}

// Allow checks if a connection from the given IP is allowed
func (crl *ConnectionRateLimiter) Allow(ip string) bool {
	crl.mutex.RLock()
	limiter, exists := crl.limiters[ip]
	crl.mutex.RUnlock()

	if !exists {
		crl.mutex.Lock()
		// Double-check after acquiring write lock
		if limiter, exists = crl.limiters[ip]; !exists {
			limiter = NewRateLimiter(crl.maxTokens, crl.refillRate)
			crl.limiters[ip] = limiter
		}
		crl.mutex.Unlock()
	}

	return limiter.Allow()
}

// Cleanup removes old rate limiters to prevent memory leaks
func (crl *ConnectionRateLimiter) Cleanup() {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	now := time.Now()
	for ip, limiter := range crl.limiters {
		limiter.mutex.Lock()
		// Remove limiters that haven't been used in the last hour
		if now.Sub(limiter.lastRefill) > time.Hour {
			delete(crl.limiters, ip)
		}
		limiter.mutex.Unlock()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
