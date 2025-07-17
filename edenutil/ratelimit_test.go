package edenutil

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	// Create a rate limiter with 2 tokens, refill every 100ms
	rl := NewRateLimiter(2, 100*time.Millisecond)

	// Should allow first two requests
	if !rl.Allow() {
		t.Error("First request should be allowed")
	}
	if !rl.Allow() {
		t.Error("Second request should be allowed")
	}

	// Third request should be denied
	if rl.Allow() {
		t.Error("Third request should be denied")
	}

	// Wait for refill and try again
	time.Sleep(150 * time.Millisecond)
	if !rl.Allow() {
		t.Error("Request after refill should be allowed")
	}
}

func TestConnectionRateLimiter_Allow(t *testing.T) {
	// Create a connection rate limiter with 2 tokens per IP, refill every 100ms
	crl := NewConnectionRateLimiter(2, 100*time.Millisecond)

	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	// Should allow first two requests from IP1
	if !crl.Allow(ip1) {
		t.Error("First request from IP1 should be allowed")
	}
	if !crl.Allow(ip1) {
		t.Error("Second request from IP1 should be allowed")
	}

	// Third request from IP1 should be denied
	if crl.Allow(ip1) {
		t.Error("Third request from IP1 should be denied")
	}

	// But requests from IP2 should still be allowed
	if !crl.Allow(ip2) {
		t.Error("First request from IP2 should be allowed")
	}
	if !crl.Allow(ip2) {
		t.Error("Second request from IP2 should be allowed")
	}

	// Wait for refill
	time.Sleep(150 * time.Millisecond)
	if !crl.Allow(ip1) {
		t.Error("Request from IP1 after refill should be allowed")
	}
}

func TestConnectionRateLimiter_Cleanup(t *testing.T) {
	crl := NewConnectionRateLimiter(1, time.Millisecond)

	// Add some IPs
	crl.Allow("192.168.1.1")
	crl.Allow("192.168.1.2")

	// Check that limiters exist
	if len(crl.limiters) != 2 {
		t.Errorf("Expected 2 limiters, got %d", len(crl.limiters))
	}

	// Wait for cleanup time and run cleanup
	time.Sleep(2 * time.Millisecond)
	crl.Cleanup()

	// Limiters should still exist since cleanup threshold is 1 hour
	if len(crl.limiters) != 2 {
		t.Errorf("Expected 2 limiters after cleanup, got %d", len(crl.limiters))
	}
}
