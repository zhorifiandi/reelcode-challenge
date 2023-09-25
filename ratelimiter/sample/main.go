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
		userID := r.URL.Query().Get("userID")
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
		time.Sleep(500 * time.Millisecond)
		w.Write([]byte("Successfully processed the message"))
	})
}
