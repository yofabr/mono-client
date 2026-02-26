package middleware

import (
	"testing"
	"time"
)

func TestRateLimiterAllowWithinWindow(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)
	now := time.Now()

	if !rl.Allow("1.1.1.1", now) {
		t.Fatalf("expected first request to be allowed")
	}

	if !rl.Allow("1.1.1.1", now.Add(time.Second)) {
		t.Fatalf("expected second request to be allowed")
	}

	if rl.Allow("1.1.1.1", now.Add(2*time.Second)) {
		t.Fatalf("expected third request in same window to be blocked")
	}
}

func TestRateLimiterResetsAfterWindow(t *testing.T) {
	rl := NewRateLimiter(1, time.Second)
	now := time.Now()

	if !rl.Allow("1.1.1.1", now) {
		t.Fatalf("expected first request to be allowed")
	}

	if rl.Allow("1.1.1.1", now.Add(500*time.Millisecond)) {
		t.Fatalf("expected second request in same window to be blocked")
	}

	if !rl.Allow("1.1.1.1", now.Add(2*time.Second)) {
		t.Fatalf("expected request to be allowed in a new window")
	}
}
