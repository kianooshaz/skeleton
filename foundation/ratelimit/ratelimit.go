package ratelimit

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter is the main struct for the sliding window rate limiter.
type RateLimiter struct {
	redisClient redis.Cmdable
	key         string
	limit       int
	window      time.Duration
	ttl         time.Duration
}

var slidingWindowLua = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])

-- Remove old entries
redis.call("ZREMRANGEBYSCORE", key, 0, now - window)

-- Count current entries
local count = redis.call("ZCARD", key)

if count < limit then
    -- Add new entry
    redis.call("ZADD", key, now, tostring(now))
    -- Set expire for safety
    redis.call("PEXPIRE", key, window)
    return 1
else
    return 0
end
`

type RateLimiterConfig struct {
	Key    string        `yaml:"key" validate:"required"`
	Limit  int           `yaml:"limit" validate:"required"`
	Window time.Duration `yaml:"window" validate:"required"`
	TTL    time.Duration `yaml:"ttl" validate:"required"`
}

func (c RateLimiterConfig) Validate() error {
	if c.Key == "" {
		return errors.New("key must be provided")
	}
	if c.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}
	if c.Window <= 0 {
		return errors.New("window must be greater than 0")
	}
	if c.TTL <= 0 {
		return errors.New("ttl must be greater than 0")
	}
	return nil
}

// New creates a new Sliding Window RateLimiter instance.
func New(redisClient redis.Cmdable, config RateLimiterConfig) (*RateLimiter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &RateLimiter{
		redisClient: redisClient,
		key:         config.Key,
		limit:       config.Limit,
		window:      config.Window,
		ttl:         config.TTL,
	}, nil
}

// Allow checks if a new request is allowed under the current sliding window rate limit.
// Returns true if allowed, false otherwise.
func (rl *RateLimiter) Allow(ctx context.Context) (bool, error) {
	nowMs := time.Now().UnixMilli()

	result, err := rl.redisClient.Eval(ctx, slidingWindowLua, []string{rl.key},
		nowMs,
		rl.window.Milliseconds(),
		rl.limit,
	).Result()

	if err != nil {
		return false, err
	}

	allowed, ok := result.(int64)
	if !ok {
		return false, nil
	}

	return allowed == 1, nil
}
