package httpclient

import (
	"net/http"
	"sync"
)

type header struct {
	mu       sync.RWMutex
	metadata http.Header
}

var _ Header = (*header)(nil)

func (h *header) Add(key string, value string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.metadata.Add(key, value)
}

func (h *header) Set(key string, value string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.metadata.Set(key, value)
}

func (h *header) Values() map[string][]string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.metadata.Clone()
}

func (h *header) Delete(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.metadata.Del(key)
}

func (h *header) Get(key string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	val := h.metadata.Get(key)
	if val == "" {
		return val, false
	}

	return val, true
}
