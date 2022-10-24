package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/exp/maps"
)

type (
	RequestOption func(*request)

	request struct {
		url     *url.URL
		payload io.Reader
		method  string
		header  http.Header
	}
)

var _ Request = (*request)(nil)

func WithPayload(payload io.Reader) RequestOption {
	return func(r *request) {
		r.payload = payload
	}
}

func WithQueryParam(param, value string) RequestOption {
	return func(r *request) {
		q := r.URL().Query()
		q.Add(param, value)
		r.url.RawQuery = q.Encode()
	}
}

func WithHeader(header map[string][]string) RequestOption {
	return func(r *request) {
		for k, val := range header {
			for _, v := range val {
				r.header.Add(k, v)
			}
		}
	}
}

func NewRequest(rawURL string, method string, opts ...RequestOption) (Request, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("[ulr.Parse] returned error: %+w", err)
	}

	r := &request{
		url:    url,
		method: method,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

func (r *request) URL() *url.URL {
	return r.url
}

func (r *request) Payload() io.Reader {
	return r.payload
}

func (r *request) Method() string {
	return r.method
}

func (r *request) Header() map[string][]string {
	snapshot := make(map[string][]string, len(r.header))
	maps.Copy(r.header, snapshot)
	return snapshot
}
