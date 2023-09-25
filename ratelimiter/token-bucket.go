package ratelimiter

import (
	"errors"
	"sync"
)

type tokenBucketRateLimiter struct {
	lock            sync.Mutex
	Limit           int
	RequestCounters map[string]int
}

func NewTokenBucketRateLimiter(limit int) *tokenBucketRateLimiter {
	rl := &tokenBucketRateLimiter{
		Limit:           limit,
		RequestCounters: map[string]int{},
	}

	return rl
}

func (rl *tokenBucketRateLimiter) getCurrentCounter(key string) int {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	counter, ok := rl.RequestCounters[key]
	if !ok {
		return 0
	}

	return counter
}

func (rl *tokenBucketRateLimiter) incrementCounter(key string) error {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	_, ok := rl.RequestCounters[key]
	if !ok {
		rl.RequestCounters[key] = 0
	}

	// log.Printf("incrementing counter from %d to %d\n", counter, counter+1)
	rl.RequestCounters[key]++

	return nil
}

func (rl *tokenBucketRateLimiter) Request(key string) (bool, error) {
	if rl.getCurrentCounter(key) < rl.Limit {
		rl.incrementCounter(key)
		return true, nil
	}

	return false, nil
}

func (rl *tokenBucketRateLimiter) Release(key string) error {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	_, ok := rl.RequestCounters[key]
	if !ok {
		return errors.New("key not found")
	}

	// log.Printf("decrementing counter from %d to %d\n", counter, counter-1)
	rl.RequestCounters[key]--

	return nil
}
