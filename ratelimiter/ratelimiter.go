package ratelimiter

type RateLimiter interface {
	Request() (bool, error)
	Release()
}
