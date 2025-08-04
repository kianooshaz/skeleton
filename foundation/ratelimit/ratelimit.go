package ratelimit

import (
	"context"
	"log/slog"
	"time"

	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/redis/go-redis/v9"
)

// RateLimiter is the main struct for the sliding window rate limiter.
type RateLimiter struct {
	redisClient redis.Cmdable
	limit       int
	window      time.Duration
	ttl         time.Duration
}

var SlidingWindowRateLimiter *RateLimiter

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
	Limit  int           `yaml:"limit" validate:"required"`
	Window time.Duration `yaml:"window" validate:"required"`
	TTL    time.Duration `yaml:"ttl" validate:"required"`
}

// NewRateLimiter creates a new Sliding Window RateLimiter instance using dependency injection.
// The redisClient parameter is the Redis client to use for rate limiting.
// The cfg parameter contains the rate limiter configuration.
// Returns a configured RateLimiter instance.
func NewRateLimiter(redisClient redis.Cmdable, cfg RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		limit:       cfg.Limit,
		window:      cfg.Window,
		ttl:         cfg.TTL,
	}
}

// NewRateLimiterFromConfig creates a new RateLimiter instance by loading configuration.
// The redisClient parameter is the Redis client to use.
// The configLoader parameter is a function that loads the rate limiter configuration.
// Returns a configured RateLimiter instance or an error.
func NewRateLimiterFromConfig(redisClient redis.Cmdable, configLoader func() (RateLimiterConfig, error)) (*RateLimiter, error) {
	cfg, err := configLoader()
	if err != nil {
		return nil, err
	}
	return NewRateLimiter(redisClient, cfg), nil
}

// Init creates a new Sliding Window RateLimiter instance and sets it globally.
// Deprecated: Use NewRateLimiter or NewRateLimiterFromConfig with dependency injection instead.
func Init(redisClient redis.Cmdable) {
	cfg, err := config.Load[RateLimiterConfig]("ratelimit")
	if err != nil {
		slog.Error("failed to load rate limiter config", "error", err)
		return
	}

	SlidingWindowRateLimiter = &RateLimiter{
		redisClient: redisClient,
		limit:       cfg.Limit,
		window:      cfg.Window,
		ttl:         cfg.TTL,
	}
}

// Allow checks if a new request is allowed under the current sliding window rate limit.
// Returns true if allowed, false otherwise.
func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	nowMs := time.Now().UnixMilli()

	result, err := rl.redisClient.Eval(ctx, slidingWindowLua, []string{key},
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
