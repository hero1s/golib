package limit

import (
	"github.com/gin-gonic/gin"
	"github.com/hero1s/golib/web/util"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter .
func NewIPRateLimiter(CountPerSecond int) *IPRateLimiter {
	var r rate.Limit
	r = 1
	b := CountPerSecond
	IpRateLimiter := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
	return IpRateLimiter
}

// AddIP creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddIP to add IP address to the map
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]
	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}
	i.mu.Unlock()
	return limiter
}

func GinIpLimit(CountPerSecond int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(CountPerSecond)
	return func(c *gin.Context) {
		ipAddr := util.GetRealIp(c)
		limiter := limiter.GetLimiter(ipAddr)
		if !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
		} else {
			c.Next()
		}
	}
}
