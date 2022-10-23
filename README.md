
# HTTP Client

It is a simple and powerful client. My motivation for start this project was that the default implementation of the NET/HTTP is highly verbose and doesn't support middleware.

With this, you will be capable to use features. For example, an error handler that is based on an HTTP status code, middleware that can be applied before and after HTTP request, and, a more straightforward public interface with total compatibility with the NET/HTTP.

## Usage/Examples

### Simple Request with Body Bytes
This implementation abstract the need to read the payload returned from the response. 

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aawadallak/httpclient"
)

const API_URL = "https://jsonplaceholder.typicode.com/todos/1"

func main() {
	ctx := context.Background()

	client := httpclient.NewClient()

	req, err := httpclient.NewRequest(API_URL, http.MethodGet)
	if err != nil {
		log.Fatalf("[httpclient.NewRequest] returned error: %+v", err)
	}

	res, err := client.Fetch(ctx, req)
	if err != nil {
		log.Fatalf("[client.Fetch] returned error: %+v", err)
	}

	body, err := res.Body().Bytes()
	if err != nil {
		log.Fatalf("[res.Body().Bytes] returned error: %+v", err)
	}

	log.Println(res.ContentLength())
	log.Println(res.Header().Values())
	log.Println(res.StatusCode())
	log.Println(string(body))
}
```

### Simple Request with Body Raw
It's the default Go body response. You **MUST** close the response body, 
otherwise, the connection will never be released from the connection pool. 

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aawadallak/httpclient"
)

const API_URL = "https://jsonplaceholder.typicode.com/todos/1"

func main() {
	ctx := context.Background()

	client := httpclient.NewClient()

	req, err := httpclient.NewRequest(API_URL, http.MethodGet)
	if err != nil {
		log.Fatalf("[httpclient.NewRequest] returned error: %+v", err)
	}

	res, err := client.Fetch(ctx, req)
	if err != nil {
		log.Fatalf("[client.Fetch] returned error: %+v", err)
	}

	// The body should be closed, otherwise,
	//the connection will never be released from net/http pool.
	defer res.Body().Close()

	resBody := res.Body()
	body, err := io.ReadAll(resBody.Raw())
	if err != nil {
		log.Fatalf("[io.ReadAll] returned error: %+v", err)
	}

	log.Println(res.ContentLength())
	log.Println(res.Header().Values())
	log.Println(res.StatusCode())
	log.Println(string(body))
}
```

### Error Handler
It is possible to set up functions to be called based on the HTTP status code. 
It's attached to the client using the functional options on the initialization. 

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aawadallak/httpclient"
)

const API_URL = "https://jsonplaceholder.typicode.com/todos/1"

func MustReturnError() httpclient.ErrorHandler {
	return func(response httpclient.Response) error {
		return errors.New("must return error")
	}
}

func main() {
	ctx := context.Background()

	client := httpclient.NewClient(
		httpclient.WithErrorHandler(http.StatusOK, MustReturnError()),
	)

	req, err := httpclient.NewRequest(
		API_URL, http.MethodGet, httpclient.WithQueryParam("foo", "bar"))
	if err != nil {
		log.Fatalf("[httpclient.NewRequest] returned error: %+v", err)
	}

	res, err := client.Fetch(ctx, req)
	if err != nil {
		log.Fatalf("[client.Fetch] returned error: %+v", err)
	}

	// The body should be closed, otherwise,
	//the connection will never be released from net/http pool.
	defer res.Body().Close()

	body, err := io.ReadAll(res.Body().Raw())
	if err != nil {
		log.Fatalf("[io.ReadAll] returned error: %+v", err)
	}

	log.Println(res.ContentLength())
	log.Println(res.Header().Values())
	log.Println(res.StatusCode())
	log.Println(string(body))
}
```

### Middlewares
Middlewares are handlers that are used to execute actions before or after the request, 
it's a powerful feature that can be used for logging, telemetry and more. 

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aawadallak/httpclient"
)

const API_URL = "https://jsonplaceholder.typicode.com/todos/1"

func MiddlewareExample() httpclient.MiddlewareHandlerFunc {
	return func(next httpclient.MiddlewareHandler) httpclient.MiddlewareHandler {
		return func(ctx context.Context, req httpclient.Request) (httpclient.Response, error) {
			log.Println("Before HTTP Request")

			res, err := next(ctx, req)
			if err != nil {
				return nil, err
			}

			res.Header().Add("X-Middleware", "middleware-example")

			log.Println("After HTTP Request")

			return res, nil
		}
	}
}

func main() {
	ctx := context.Background()

	client := httpclient.NewClient(
		httpclient.WithMiddleware(MiddlewareExample()),
	)

	req, err := httpclient.NewRequest(API_URL, http.MethodGet)
	if err != nil {
		log.Fatalf("[httpclient.NewRequest] returned error: %+v", err)
	}

	res, err := client.Fetch(ctx, req)
	if err != nil {
		log.Fatalf("[client.Fetch] returned error: %+v", err)
	}

	// The body should be closed, otherwise,
	//the connection will never be released from net/http pool.
	defer res.Body().Close()

	body, err := io.ReadAll(res.Body().Raw())
	if err != nil {
		log.Fatalf("[io.ReadAll] returned error: %+v", err)
	}

	log.Println(res.ContentLength())
	log.Println(res.Header().Values())
	log.Println(res.StatusCode())

	log.Println(string(body))
}
```
## Authors

- [@aawadallak](https://www.github.com/aawadallak)

