package ratelimiter_test

import (
	"fmt"
	"log"
	"net/http"
	"ratelimiter"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketRateLimiter(t *testing.T) {
	LIMIT := 5

	rl := ratelimiter.NewTokenBucketRateLimiter(LIMIT)
	go sampleWebServer(rl)
	time.Sleep(1 * time.Second)

	// Try sending 10 requests to the server
	// The first 5 requests should be accepted
	// The next 5 requests should be rejected
	TOTAL_REQUESTS := 10
	resultChan := make(chan bool, TOTAL_REQUESTS)
	for i := 0; i < TOTAL_REQUESTS; i++ {
		go func(resultChan *chan bool, requestID int) {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", SAMPLE_WEB_SERVER_PORT))
			if err != nil {
				t.Errorf("Error: %+v\n", err)
			}

			if resp.StatusCode == http.StatusOK {
				log.Printf("Request %d was accepted\n", requestID+1)
				*resultChan <- true
			} else {
				log.Printf("Request %d was rejected\n", requestID+1)
				*resultChan <- false
			}
		}(&resultChan, i)
	}

	// Wait for all the requests to complete
	succeeded := 0
	failed := 0
	for i := 0; i < TOTAL_REQUESTS; i++ {
		result := <-resultChan
		if result {
			succeeded++
		} else {
			failed++
		}
	}

	assert.Equal(t, succeeded, LIMIT)
	assert.Equal(t, failed, TOTAL_REQUESTS-LIMIT)
}
