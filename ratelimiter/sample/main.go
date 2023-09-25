package main

import (
	"log"
	"net/http"
	"ratelimiter"
	"time"
)

func main() {
	rl := ratelimiter.NewTokenBucketRateLimiter(5)

	router := http.NewServeMux()
	router.Handle("/", index(rl))

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

func index(rl ratelimiter.RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		time.Sleep(5 * time.Second)
		w.Write([]byte("Successfully processed the message"))
	})
}
