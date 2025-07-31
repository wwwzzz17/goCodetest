package services

import (
	"sync"
	"time"
)

type RateLimiter struct {
	limit   int
	window  time.Duration
	allowed []time.Time
	mu      sync.Mutex
}

func NewRateLimiter(n int) *RateLimiter {
	return &RateLimiter{
		limit:   n,
		window:  time.Second,
		allowed: make([]time.Time, 0),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	startedTime := now.Add(-rl.window)

	inTimeWindow := 0
	for i := len(rl.allowed) - 1; i >= 0; i-- {
		if rl.allowed[i].After(startedTime) {
			inTimeWindow++
		} else {
			break
		}
	}

	if inTimeWindow < len(rl.allowed) {
		rl.allowed = rl.allowed[len(rl.allowed)-inTimeWindow:]
	}

	if len(rl.allowed) >= rl.limit {
		return false
	}

	rl.allowed = append(rl.allowed, now)
	return true
}
