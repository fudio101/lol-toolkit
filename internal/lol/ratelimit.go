package lol

import (
	"sync"
	"time"
)

// Rate limit constants from Riot API.
const (
	shortWindow = time.Second
	shortLimit  = 20 // 20 requests per second

	longWindow = 2 * time.Minute
	longLimit  = 100 // 100 requests per 2 minutes
)

// RateLimiter implements Riot API rate limiting with sliding windows.
type RateLimiter struct {
	mu            sync.Mutex
	shortRequests []time.Time
	longRequests  []time.Time
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		shortRequests: make([]time.Time, 0, shortLimit),
		longRequests:  make([]time.Time, 0, longLimit),
	}
}

// Wait blocks until a request can be made within rate limits.
func (r *RateLimiter) Wait() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	r.cleanup(now)

	// Wait for short limit if needed
	if len(r.shortRequests) >= shortLimit {
		r.waitForSlot(r.shortRequests[0], shortWindow, &now)
	}

	// Wait for long limit if needed
	if len(r.longRequests) >= longLimit {
		r.waitForSlot(r.longRequests[0], longWindow, &now)
	}

	// Record request
	r.shortRequests = append(r.shortRequests, now)
	r.longRequests = append(r.longRequests, now)
}

// waitForSlot waits until a slot is available.
func (r *RateLimiter) waitForSlot(oldest time.Time, window time.Duration, now *time.Time) {
	waitTime := window - now.Sub(oldest)
	if waitTime <= 0 {
		return
	}

	r.mu.Unlock()
	time.Sleep(waitTime)
	r.mu.Lock()

	*now = time.Now()
	r.cleanup(*now)
}

// cleanup removes expired requests from both windows.
func (r *RateLimiter) cleanup(now time.Time) {
	r.shortRequests = filterAfter(r.shortRequests, now.Add(-shortWindow))
	r.longRequests = filterAfter(r.longRequests, now.Add(-longWindow))
}

// filterAfter returns only timestamps after cutoff.
func filterAfter(times []time.Time, cutoff time.Time) []time.Time {
	result := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			result = append(result, t)
		}
	}
	return result
}

// CanMakeRequest returns true if a request can be made immediately.
func (r *RateLimiter) CanMakeRequest() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cleanup(time.Now())
	return len(r.shortRequests) < shortLimit && len(r.longRequests) < longLimit
}

// GetStatus returns current rate limit usage.
func (r *RateLimiter) GetStatus() (shortUsed, shortMax, longUsed, longMax int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cleanup(time.Now())
	return len(r.shortRequests), shortLimit, len(r.longRequests), longLimit
}
