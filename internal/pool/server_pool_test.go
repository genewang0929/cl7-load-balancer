package pool

import (
	"strconv"
	"testing"
)

// 測試 Round-Robin 演算法是否正確輪詢
func TestGetNextPeer_RoundRobin(t *testing.T) {
	pool := &ServerPool{}
	pool.AddBackend("http://localhost:8081")
	pool.AddBackend("http://localhost:8082")
	pool.AddBackend("http://localhost:8083")

	// 預期順序應該是 8081 -> 8082 -> 8083 -> 8081
	expectedPorts := []string{"8081", "8082", "8083", "8081"}

	for i, port := range expectedPorts {
		peer := pool.GetNextPeer()
		if peer == nil {
			t.Fatalf("第 %d 次嘗試時 peer 不應為 nil", i)
		}
		if peer.URL.Port() != port {
			t.Errorf("第 %d 次嘗試預期 Port 為 %s，但得到 %s", i, port, peer.URL.Port())
		}
	}
}

// 測試當部分後端失效時，是否會自動跳過
func TestGetNextPeer_SkipDeadServer(t *testing.T) {
	pool := &ServerPool{}
	pool.AddBackend("http://localhost:8081")
	pool.AddBackend("http://localhost:8082")
	pool.AddBackend("http://localhost:8083")

	// 模擬 8082 壞掉了
	pool.Backends[1].SetAlive(false)

	// 預期順序應該跳過 8082：8081 -> 8083 -> 8081
	expectedPorts := []string{"8081", "8083", "8081"}

	for _, port := range expectedPorts {
		peer := pool.GetNextPeer()
		if peer.URL.Port() != port {
			t.Errorf("跳過失效伺服器測試失敗：預期 %s，得到 %s", port, peer.URL.Port())
		}
	}
}

// 測試當所有伺服器都掛掉時，是否回傳 nil
func TestGetNextPeer_AllDown(t *testing.T) {
	pool := &ServerPool{}
	pool.AddBackend("http://localhost:8081")
	pool.Backends[0].SetAlive(false)

	peer := pool.GetNextPeer()
	if peer != nil {
		t.Error("當所有後端都失效時，應該回傳 nil")
	}
}

// 效能測試：評估高併發下 GetNextPeer 的表現
func BenchmarkGetNextPeer(b *testing.B) {
	pool := &ServerPool{}
	// 模擬 10 台後端伺服器
	for i := 0; i < 10; i++ {
		pool.AddBackend("http://localhost:" + strconv.Itoa(8080+i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.GetNextPeer()
		}
	})
}
