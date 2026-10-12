package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*Client)
)

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]

	if !exists {
		limiter := rate.NewLimiter(5, 10) // 5 req/sec, bursts = 10
		newClient := &Client{limiter, time.Now()}
		clients[ip] = newClient
		// log.Printf("A clients[%s] - {limiter: %+v, lastSeen: %s}", ip, newClient.limiter, newClient.lastSeen)

		return limiter
	}

	// log.Printf("A clients[%s] - {limiter: %+v, lastSeen: %s}", ip, client.limiter, client.lastSeen)
	client.lastSeen = time.Now()
	return client.limiter
}

func CleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// ab -n 20 -c 1 -H 'X-API-Key: 5ae486dd-1026-4ee0-9ce2-be601e25d64b' http://localhost:8080/api/v1/categories/golang

func RateLimitingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)
		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Vui lòng thử lại sau",
			})

			return
		}

		ctx.Next()
	}
}
