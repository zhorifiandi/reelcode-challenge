package ratelimiter_test

import (
	"fmt"
	"log"
	"net/http"
	"ratelimiter"
	"time"
)

const SAMPLE_WEB_SERVER_PORT = 9990

func sampleWebServer(rl ratelimiter.RateLimiter) error {
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.URL.Query().Get("requestID")
		userID := r.URL.Query().Get("userID")
		log.Printf("Incoming request on this route: %+v, for this user: %+v, requestID: %+v\n", r.URL, userID, requestID)

		ok, err := rl.Request(userID)
		if err != nil {
			log.Printf("Error: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !ok {
			log.Printf("Request rejected on this route: %+v, for this user: %+v\n", r.URL, userID)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		defer rl.Release(userID)
		log.Printf("Request accepted on this route: %+v, for this user: %+v\n", r.URL, userID)
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("Successfully processed the message"))
	}))

	err := http.ListenAndServe(fmt.Sprintf(":%d", SAMPLE_WEB_SERVER_PORT), router)
	if err != nil {
		panic(err)
	}
	return err
}

func makeRequestFunc(resultChan *chan bool, requestID int, userID string) error {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d?userID=%s&requestID=%d", SAMPLE_WEB_SERVER_PORT, userID, requestID))
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		log.Printf("Request %d was accepted\n", requestID+1)
		*resultChan <- true
	} else {
		log.Printf("Request %d was rejected\n", requestID+1)
		*resultChan <- false
	}

	return nil
}
