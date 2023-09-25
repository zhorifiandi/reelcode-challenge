package ratelimiter_test

import (
	"fmt"
	"log"
	"net/http"
	"ratelimiter"
	"time"
)

const SAMPLE_WEB_SERVER_PORT = 9999

func sampleWebServer(rl ratelimiter.RateLimiter) error {
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok, err := rl.Request()
		if err != nil {
			log.Printf("Error: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !ok {
			log.Printf("Request rejected on this route: %+v\n", r.URL)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		defer rl.Release()
		log.Printf("Request accepted on this route: %+v\n", r.URL)
		time.Sleep(500 * time.Millisecond)
		w.Write([]byte("Successfully processed the message"))
	}))

	err := http.ListenAndServe(fmt.Sprintf(":%d", SAMPLE_WEB_SERVER_PORT), router)
	if err != nil {
		panic(err)
	}
	return err
}
