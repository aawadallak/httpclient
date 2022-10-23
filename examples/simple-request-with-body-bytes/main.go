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

	fmt.Println(res.ContentLength())
	fmt.Println(res.Header().Values())
	fmt.Println(res.StatusCode())
	fmt.Println(string(body))
}
