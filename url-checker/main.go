package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"sync"
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

func printResult(result Result) {
	if result.Error != nil {
		fmt.Printf(
			"%s - ERROR: %v\n",
			result.URL,
			result.Error,
		)
		return
	}

	fmt.Printf(
		"%s - HTTP %d\n",
		result.URL,
		result.Status,
	)
}

func main() {
	jobs := make(chan string)
	results := make(chan Result)
		
	var wg sync.WaitGroup

	ctx := context.Background()

	/*

	jobs <- "https://www.google.com"

	urls := <-jobs

	ctx := context.Background()

	result := checkURL(ctx, urls)

	fmt.Printf("URL: %s, Status: %s, Error: %v\n", result.URL, result.Status, result.Error)

	*/
	for range 8 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for url := range jobs {
				result := checkURL(ctx, url)
				results <- result
			}
		}()
	}
	
	urls := []string{
		"https://www.google.com",
		"https://www.example.com",
		"https://github.com",
		"https://golang.org",
		"https://stackoverflow.com",
		"https://www.nonexistentwebsite.com",
	}

	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()	

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		printResult(result)
	}

}
