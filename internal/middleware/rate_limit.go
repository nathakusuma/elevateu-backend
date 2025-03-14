package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"

	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
)

type IpRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	r      rate.Limit
	b      int
	ttl    time.Duration
	expire map[string]time.Time
}

func NewIpRateLimiter(r rate.Limit, b int, ttl time.Duration) *IpRateLimiter {
	i := &IpRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		r:      r,
		b:      b,
		ttl:    ttl,
		expire: make(map[string]time.Time),
	}

	go i.cleanup(ttl / 2)

	return i
}

func (i *IpRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
		i.expire[ip] = time.Now().Add(i.ttl)
	} else {
		i.expire[ip] = time.Now().Add(i.ttl)
	}

	return limiter
}

func (i *IpRateLimiter) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()
		now := time.Now()
		for ip, exp := range i.expire {
			if exp.Before(now) {
				delete(i.ips, ip)
				delete(i.expire, ip)
			}
		}
		i.mu.Unlock()
	}
}

type RateLimitConfig struct {
	MaxRequests    int           // Number of requests allowed
	PerTimeWindow  time.Duration // Time window for the allowed requests
	ExpirationTime time.Duration // Time after which the limiter is removed if not used
}

func (m *Middleware) RateLimit(config RateLimitConfig) fiber.Handler {
	requestsPerSecond := float64(config.MaxRequests) / config.PerTimeWindow.Seconds()

	burstSize := config.MaxRequests

	limiter := NewIpRateLimiter(rate.Limit(requestsPerSecond), burstSize, config.ExpirationTime)

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		if ip == "" {
			ip = "unknown"
		}

		ipLimiter := limiter.GetLimiter(ip)

		if !ipLimiter.Allow() {
			return errorpkg.ErrRateLimitExceeded()
		}

		return c.Next()
	}
}
