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
	fmt.Println(string(body))
}
