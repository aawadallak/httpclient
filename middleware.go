package httpclient

func chainMiddlewares(middlewares ...MiddlewareHandlerFunc) MiddlewareHandlerFunc {
	return func(next MiddlewareHandler) MiddlewareHandler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}
}
