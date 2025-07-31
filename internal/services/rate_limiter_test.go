package services

import (
	"testing"
	"time"
)

func TestRatelimiter(t *testing.T) {
	rl := NewRateLimiter(5)

	for i := 0; i < 10; i++ {
		if i < 5 {
			if !rl.Allow() {
				t.Errorf("Expected request %d to be allowed", i)
			}
		} else {
			if rl.Allow() {
				t.Errorf("Expected request %d to be denied", i)
			}
		}
	}
}

func TestRatelimiterSleepHalfSecond(t *testing.T) {
	rl := NewRateLimiter(5)

	for i := 0; i < 10; i++ {
		time.Sleep(200 * time.Millisecond)
		if !rl.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}
}
