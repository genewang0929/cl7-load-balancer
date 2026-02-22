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
			log.Println("--- Health Check Report ---")
			for _, b := range s.Backends {
				alive := isAlive(b.URL.Host)
				b.SetAlive(alive)

				status := "up"
				if !alive {
					status = "down"
				}

				// 修改：輸出日誌時同時顯示該節點處理的請求數
				log.Printf("[%s] Status: %s, Total Requests: %d", b.URL, status, b.GetRequests())
			}
			log.Println("---------------------------")
		case <-ctx.Done():
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
