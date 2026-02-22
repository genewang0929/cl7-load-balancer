package pool

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
	Alive        bool
	mux          sync.RWMutex
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

type ServerPool struct {
	Backends []*Backend
	current  int
	mu       sync.Mutex
}

func (s *ServerPool) AddBackend(backendUrl string) {
	u, _ := url.Parse(backendUrl)
	proxy := httputil.NewSingleHostReverseProxy(u)

	// 修改 Director 以確保 Host Header 正確
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = u.Host
	}

	s.Backends = append(s.Backends, &Backend{
		URL:          u,
		ReverseProxy: proxy,
		Alive:        true,
	})
}

func (s *ServerPool) GetNextPeer() *Backend {
	s.mu.Lock()
	defer s.mu.Unlock()

	total := len(s.Backends)
	if total == 0 {
		return nil
	}

	for i := 0; i < total; i++ {
		idx := (s.current + i) % total
		if s.Backends[idx].IsAlive() {
			if i != total-1 {
				s.current = (idx + 1) % total
			}
			return s.Backends[idx]
		}
	}
	return nil
}
