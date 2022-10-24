package httpclient

import (
	"context"
	"net/http"
	"time"
)

type (
	CallOption func(*callConfig)

	callConfig struct {
		timeout time.Duration
	}
)

// WithTimeout allow to set up a duration to be used on the call
// instead of the timeout that has been defined on the initialization.
func WithTimeout(timeout time.Duration) CallOption {
	return func(c *callConfig) {
		c.timeout = timeout
	}
}

func newCallConfig(opts ...CallOption) *callConfig {
	c := &callConfig{
		timeout: time.Second * 10,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type httpClient struct {
	client *http.Client
}

func defaultHTTPClient() *httpClient {
	client := &http.Client{
		Timeout:   time.Second * 60,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}

	return &httpClient{
		client: client,
	}
}

func (h *httpClient) Do(ctx context.Context, req Request, opts ...CallOption) (Response, error) {
	h.applyRequestOptions(opts...)

	resp, err := h.call(ctx, req)
	if err != nil {
		return nil, err
	}

	return &response{
		contentLenght: resp.ContentLength,
		status:        resp.StatusCode,
		body:          &body{raw: resp.Body},
		header:        &header{metadata: req.Header()},
	}, nil
}

func (h *httpClient) applyRequestOptions(opts ...CallOption) {
	if len(opts) != 0 {
		cfg := newCallConfig(opts...)

		h.client.Timeout = cfg.timeout
	}
}

func (h *httpClient) call(ctx context.Context, req Request) (*http.Response, error) {
	r, err := http.NewRequestWithContext(
		ctx, req.Method(), req.URL().String(), req.Payload())
	if err != nil {
		return nil, err
	}

	r.Header = req.Header()

	resp, err := h.client.Do(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
