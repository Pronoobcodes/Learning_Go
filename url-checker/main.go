package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

/*

func main() {
	fmt.Println("URL Checker")
}

*/	

type Result struct {
	URL    string
	Status string
	Error  error
}

func checkURL(ctx context.Context, url string) Result {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{URL: url, Status: "", Error: err}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{URL: url, Status: "", Error: err}
	}
	defer resp.Body.Close()

	return Result{URL: url, Status: fmt.Sprintf("%d", resp.StatusCode), Error: nil}
}

func main() {
	ctx := context.Background()

	result := checkURL(ctx, "https://www.google.com")

	fmt.Printf("URL: %s, Status: %s, Error: %v\n", result.URL, result.Status, result.Error)
}
