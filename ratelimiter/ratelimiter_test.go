package ratelimiter_test

import (
	"fmt"
	"log"
	"net/http"
	"ratelimiter"
	"sync"
	"sync/atomic"
	"time"
)

const SAMPLE_WEB_SERVER_PORT = 9990

func sampleWebServer(
	rl ratelimiter.RateLimiter,
	mockedProcessingTime time.Duration,
) error {
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// requestID := r.URL.Query().Get("requestID")
		userID := r.URL.Query().Get("userID")
		// log.Printf("Incoming request on this route: %+v, for this user: %+v, requestID: %+v\n", r.URL, userID, requestID)

		ok, err := rl.Request(userID)
		if err != nil {
			log.Printf("Error: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !ok {
			// log.Printf("Request rejectedCounter on this route: %+v, for this user: %+v\n", r.URL, userID)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		defer rl.Release(userID)
		// log.Printf("Request accepted on this route: %+v, for this user: %+v\n", r.URL, userID)
		time.Sleep(mockedProcessingTime)
		w.Write([]byte("Successfully processed the message"))
	}))

	err := http.ListenAndServe(fmt.Sprintf(":%d", SAMPLE_WEB_SERVER_PORT), router)
	if err != nil {
		panic(err)
	}
	return err
}

func makeRequestFunc(
	wg *sync.WaitGroup,
	succeededCounter *atomic.Int32,
	rejectedCounter *atomic.Int32,
	errorCounter *atomic.Int32,
	requestID int,
	userID string,
) error {
	defer wg.Done()
	log.Printf("Sending request %d\n", requestID)
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d?userID=%s&requestID=%d", SAMPLE_WEB_SERVER_PORT, userID, requestID))
	if err != nil {
		log.Printf("Error: %+v\n", err)
		errorCounter.Add(1)
		return err
	}

	if resp.StatusCode == http.StatusOK {
		log.Printf("Request %d was accepted\n", requestID)
		succeededCounter.Add(1)
	} else {
		log.Printf("Request %d was rejectedCounter\n", requestID)
		rejectedCounter.Add(1)
	}

	return nil
}
