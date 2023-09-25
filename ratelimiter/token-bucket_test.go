package ratelimiter_test

import (
	"log"
	"ratelimiter"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketRateLimiter(t *testing.T) {
	LIMIT := 20

	rl := ratelimiter.NewTokenBucketRateLimiter(LIMIT)
	go sampleWebServer(rl)
	time.Sleep(1 * time.Second)

	// Try sending 100 requests to the server for User1
	// The first 20 requests should be accepted
	// The next 80 requests should be rejectedps
	userID1 := "user1"
	totalRequestsUser1 := 100
	user1ResultChan := make(chan bool, totalRequestsUser1)
	for i := 0; i < totalRequestsUser1; i++ {
		go makeRequestFunc(&user1ResultChan, i, userID1)
	}

	// Try sending 20 requests to the server for User2
	// All the requests should be accepted
	userID2 := "user2"
	totalRequestsUser2 := LIMIT
	user2ResultChan := make(chan bool, totalRequestsUser2)
	for i := 0; i < totalRequestsUser2; i++ {
		go makeRequestFunc(&user2ResultChan, i, userID2)
	}

	// Wait for all the requests to complete
	succeeded := 0
	rejected := 0
	for i := 0; i < totalRequestsUser1+totalRequestsUser2; i++ {
		select {
		case condition := <-user1ResultChan:
			if condition {
				succeeded++
			} else {
				rejected++
			}
		case condition := <-user2ResultChan:
			if condition {
				succeeded++
			} else {
				assert.Fail(t, "Request for user2 should not be rejected")
			}
		}
	}

	log.Printf("Succeeded: %d, Rejected: %d\n", succeeded, rejected)
	assert.Equal(t, succeeded, LIMIT+totalRequestsUser2)
	assert.Equal(t, rejected, totalRequestsUser1-LIMIT)
}
