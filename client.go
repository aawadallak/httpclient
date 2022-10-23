package httpclient

import (
	"context"
	"net/http"
)

type (
	Option func(*client)

	client struct {
		httpClient   *httpClient
		errorHandler map[int]ErrorHandler
		middlewares  []MiddlewareHandlerFunc
	}
)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *client) {
		c.httpClient.client = httpClient
	}
}

func WithMiddleware(middleware MiddlewareHandlerFunc) Option {
	return func(c *client) {
		c.middlewares = append(c.middlewares, middleware)
	}
}

func WithErrorHandler(status int, handler ErrorHandler) Option {
	return func(c *client) {
		c.errorHandler[status] = handler
	}
}

func NewClient(opts ...Option) Client {
	c := &client{
		httpClient:   defaultHTTPClient(),
		errorHandler: make(map[int]ErrorHandler),
		middlewares:  make([]MiddlewareHandlerFunc, 0),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *client) Fetch(ctx context.Context, req Request, opts ...CallOption) (Response, error) {
	res, err := c.callHandler(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	if err := c.checkErrorHandler(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *client) callHandler(ctx context.Context, req Request, opts ...CallOption) (Response, error) {
	handler := func(ctx context.Context, req Request) (Response, error) {
		return c.httpClient.Do(ctx, req, opts...)
	}

	middleware := chainMiddlewares(c.middlewares...)

	res, err := middleware(handler)(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *client) checkErrorHandler(res Response) error {
	fn, ok := c.errorHandler[res.StatusCode()]
	if ok {
		return fn(res)
	}

	return nil
}
