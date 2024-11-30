package api

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	limiters map[string]*Client
	mu       sync.RWMutex
}

func NewLimiter() *RateLimiter {
	limiter := RateLimiter{limiters: make(map[string]*Client)}

	cleanup := func() {
		for {
			time.Sleep(time.Minute)
			limiter.mu.Lock()
			for ip, client := range limiter.limiters {
				if time.Since(client.lastSeen) > 30*time.Second {
					delete(limiter.limiters, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}

	go cleanup()

	return &limiter
}

func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, ok := r.limiters[ip]

	if !ok {
		client = &Client{rate.NewLimiter(2, 5), time.Now()}
		r.limiters[ip] = client
	} else {
		client.lastSeen = time.Now()
	}

	return client.limiter.Allow()
}
