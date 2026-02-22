package health

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/genewang0929/cl7-load-balancer/internal/pool"
)

func RunHealthCheck(ctx context.Context, s *pool.ServerPool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, b := range s.Backends {
				alive := isAlive(b.URL.Host)
				b.SetAlive(alive)
				status := "up"
				if !alive {
					status = "down"
				}
				log.Printf("Health check: %s is %s", b.URL, status)
			}
		case <-ctx.Done():
			log.Println("Health check worker shutting down...")
			return
		}
	}
}

func isAlive(host string) bool {
	conn, err := net.DialTimeout("tcp", host, 2*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
