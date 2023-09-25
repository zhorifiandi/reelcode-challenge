package ratelimiter

import (
	"log"
	"sync"
)

type tokenBucketRateLimiter struct {
	lock           sync.Mutex
	Limit          int
	CurrentCounter int
}

func NewTokenBucketRateLimiter(limit int) *tokenBucketRateLimiter {
	rl := &tokenBucketRateLimiter{
		Limit:          limit,
		CurrentCounter: 0,
	}

	return rl
}

func (rl *tokenBucketRateLimiter) getCurrentCounter() int {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	return rl.CurrentCounter
}

func (rl *tokenBucketRateLimiter) incrementCounter() {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	log.Printf("incrementing counter from %d to %d\n", rl.CurrentCounter, rl.CurrentCounter+1)
	rl.CurrentCounter++
}

func (rl *tokenBucketRateLimiter) Request() (bool, error) {
	if rl.getCurrentCounter() < rl.Limit {
		rl.incrementCounter()
		return true, nil
	}

	return false, nil
}

func (rl *tokenBucketRateLimiter) Release() {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	log.Printf("decrementing counter from %d to %d\n", rl.CurrentCounter, rl.CurrentCounter-1)
	rl.CurrentCounter--
}
