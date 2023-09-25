package ratelimiter_test

import (
	"log"
	"ratelimiter"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketRateLimiter(t *testing.T) {
	LIMIT := 200

	rl := ratelimiter.NewTokenBucketRateLimiter(LIMIT)
	go sampleWebServer(
		rl,
		500*time.Millisecond, // Give enough buffer for the requests to be processed, to test the rate limiting
	)
	time.Sleep(1 * time.Second)

	var wg sync.WaitGroup

	// Try sending 1000 requests to the server for User1
	// The first 200 requests should be accepted
	// The next 800 requests should be rejectedps
	user1SucceededCounter := atomic.Int32{}
	user1RejectedCounter := atomic.Int32{}
	user1ErrorCounter := atomic.Int32{}
	userID1 := "user1"
	totalRequestsUser1 := 1000
	for i := 0; i < totalRequestsUser1; i++ {
		wg.Add(1)
		go makeRequestFunc(&wg, &user1SucceededCounter, &user1RejectedCounter, &user1ErrorCounter, i, userID1)
	}

	// Try sending 200 requests to the server for User2
	// All the requests should be accepted
	user2SucceededCounter := atomic.Int32{}
	user2RejectedCounter := atomic.Int32{}
	user2ErrorCounter := atomic.Int32{}
	userID2 := "user2"
	totalRequestsUser2 := LIMIT
	for i := 0; i < totalRequestsUser2; i++ {
		wg.Add(1)
		go makeRequestFunc(&wg, &user2SucceededCounter, &user2RejectedCounter, &user2ErrorCounter, i, userID2)
	}

	wg.Wait()

	user1Succeeded := int(user1SucceededCounter.Load())
	user1Rejected := int(user1RejectedCounter.Load())
	user1ErrorCount := int(user1ErrorCounter.Load())

	log.Printf("[User 1] Succeeded: %d, Rejected: %d, Error: %d\n", user1Succeeded, user1Rejected, user1ErrorCount)
	assert.LessOrEqual(t, user1Succeeded, LIMIT)
	assert.GreaterOrEqual(t, user1Rejected, 0)

	user2Succeeded := int(user2SucceededCounter.Load())
	user2Rejected := int(user2RejectedCounter.Load())
	user2ErrorCount := int(user2ErrorCounter.Load())

	log.Printf("[User 2] Succeeded: %d, Rejected: %d, Error: %d\n", user2Succeeded, user2Rejected, user2ErrorCount)
	assert.Equal(t, user2Succeeded, totalRequestsUser2-user2ErrorCount)
	assert.Equal(t, user2Rejected, 0)
}
