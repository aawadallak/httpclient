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
			fmt.Println("Before HTTP Request")

			res, err := next(ctx, req)
			if err != nil {
				return nil, err
			}

			res.Header().Add("X-Middleware", "middleware-example")

			fmt.Println("After HTTP Request")

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

	fmt.Println(string(body))
}
