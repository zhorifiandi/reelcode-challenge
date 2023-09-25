package ratelimiter_test

import (
	"ratelimiter"
	"testing"
)

func TestTokenBucketRateLimiter(t *testing.T) {
	LIMIT := 200

	rl := ratelimiter.NewTokenBucketRateLimiter(LIMIT)
	runRateLimiterTestCase(t, rl, LIMIT)
}
