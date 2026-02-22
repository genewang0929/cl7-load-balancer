package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/genewang0929/cl7-load-balancer/internal/health"
	"github.com/genewang0929/cl7-load-balancer/internal/pool"
)

func main() {
	serverPool := &pool.ServerPool{}
	// 這裡示範手動加入，實務上可從 config.yaml 讀取
	backends := []string{"http://localhost:8081", "http://localhost:8082"}
	for _, b := range backends {
		serverPool.AddBackend(b)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動健康檢查
	go health.RunHealthCheck(ctx, serverPool, 10*time.Second)

	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			peer := serverPool.GetNextPeer()
			if peer != nil {
				// 新增：使用原子操作增加計數，確保併發安全
				atomic.AddUint64(&peer.Requests, 1)

				peer.ReverseProxy.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Service not available", http.StatusServiceUnavailable)
		}),
	}

	// 監聽系統訊號以進行優雅關閉
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		cancel() // 停止健康檢查

		shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
	}()

	log.Println("Load Balancer started on :8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
