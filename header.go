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

// Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
// The key is case insensitive; it is canonicalized by CanonicalHeaderKey.
//
// It's concurrent safe to add keys to the header.
func (h *header) Add(key string, value string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.metadata.Add(key, value)
}

// Set sets the header entries associated with key to the single element value.
// It replaces any existing values associated with key.
// The key is case insensitive; it is canonicalized by textproto.CanonicalMIMEHeaderKey.
//
// It's concurrent safe to set keys to the header.
func (h *header) Set(key string, value string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.metadata.Set(key, value)
}

// Del deletes the values associated with key.
// The key is case insensitive; it is canonicalized by CanonicalHeaderKey.
//
// It's concurrent safe to delete keys from the header.
func (h *header) Delete(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.metadata.Del(key)
}

// Values return a copy from the header content.
//
// It's concurrent safe to get values from the header.
func (h *header) Values() map[string][]string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.metadata.Clone()
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns an empty string.
// It is case insensitive; textproto.CanonicalMIMEHeaderKey is used to canonicalize the provided key.
//
// It's concurrent safe to get a value from the header.
func (h *header) Get(key string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	val := h.metadata.Get(key)
	if val == "" {
		return val, false
	}

	return val, true
}
