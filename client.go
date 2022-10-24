package httpclient

import (
	"context"
	"net/http"
)

type (
	// Option is used to extend the client behavior.
	Option func(*client)

	client struct {
		httpClient   *httpClient
		errorHandler map[int]ErrorHandler
		middlewares  []MiddlewareHandlerFunc
	}
)

// WithHTTPClient allow a custom HTTP client to be used instead of the default client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *client) {
		c.httpClient.client = httpClient
	}
}

// WithMiddleware allow pass a middleware to be used by the high-level client before or
// after the HTTP Request.
//
// The middlewares will be called in the order they are passed to the intialization.
func WithMiddleware(middleware MiddlewareHandlerFunc) Option {
	return func(c *client) {
		c.middlewares = append(c.middlewares, middleware)
	}
}

// WithErrorHandler allow a handler to be called based on the informed HTTP status.
func WithErrorHandler(status int, handler ErrorHandler) Option {
	return func(c *client) {
		c.errorHandler[status] = handler
	}
}

// NewClient initializes a new client using the provided configuration.
// It's extensible by using funcitional options.
//
// The available options are:
//
//   - WithHTTPClient(httpClient *http.Client)
//   - WithMiddleware(middleware MiddlewareHandlerFunc)
//   - WithErrorHandler(status int, handler ErrorHandler)
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

// Fetch start a request using the giving configuration. It's extensible by using functional options.
//
// The middlewares will be applied first and then the error handling will be triggered whether
// it's configured, otherwise, the client will return the response received from the request.
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
