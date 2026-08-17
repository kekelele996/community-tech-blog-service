package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/util"
)

// tokenBucket 内存令牌桶限流器
type tokenBucket struct {
	mu       sync.Mutex
	tokens   map[string]*bucket
	capacity float64
	rate     float64 // 每秒补充令牌数
}

type bucket struct {
	tokens float64
	last   time.Time
}

func newTokenBucket(capacity, rate float64) *tokenBucket {
	return &tokenBucket{tokens: make(map[string]*bucket), capacity: capacity, rate: rate}
}

func (tb *tokenBucket) allow(key string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	now := time.Now()
	b, ok := tb.tokens[key]
	if !ok {
		tb.tokens[key] = &bucket{tokens: tb.capacity, last: now}
		return true
	}
	elapsed := now.Sub(b.last).Seconds()
	b.tokens = min(tb.capacity, b.tokens+elapsed*tb.rate)
	b.last = now
	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// RateLimit 限流中间件：按 IP 限流（默认 60 次/分钟）
func RateLimit(limit float64) gin.HandlerFunc {
	tb := newTokenBucket(limit, limit/60)
	return func(c *gin.Context) {
		ip := util.ClientIP(c)
		if !tb.allow(ip) {
			util.LogInfo(util.GetRequestID(c), "request rate limited", "ip", ip, "path", c.Request.URL.Path)
			util.Fail(c, 429, constants.CodeRateLimited, constants.MessageOf(constants.CodeRateLimited))
			return
		}
		c.Next()
	}
}
