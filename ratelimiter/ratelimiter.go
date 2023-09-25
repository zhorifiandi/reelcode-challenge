package ratelimiter

type RateLimiter interface {
	Request(key string) (bool, error)
	Release(key string) error
}
