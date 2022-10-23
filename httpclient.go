package httpclient

import (
	"context"
	"io"
	"net/url"
)

type (
	Request interface {
		URL() *url.URL
		Payload() io.Reader
		Header() map[string][]string
		Method() string
	}

	Response interface {
		Body() Body
		ContentLength() int64
		StatusCode() int
		Header() Header
	}

	Header interface {
		Add(key, value string)
		Set(key, value string)
		Values() map[string][]string
		Delete(key string)
		Get(key string) (string, bool)
	}

	Body interface {
		Raw() io.Reader
		Close() error
		Bytes() ([]byte, error)
	}

	MiddlewareHandler func(ctx context.Context, req Request) (Response, error)

	MiddlewareHandlerFunc func(next MiddlewareHandler) MiddlewareHandler

	ErrorHandler func(r Response) error

	Client interface {
		Fetch(ctx context.Context, req Request, opts ...CallOption) (Response, error)
	}
)
