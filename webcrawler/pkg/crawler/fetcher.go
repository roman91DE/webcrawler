package crawler

import (
	"fmt"
	"io"
	"net/http"
	// "regexp"
	"time"
)

type FetchResult struct {
	URL     string
	Status  int
	Content string
	Error   error
}

func FetchAll(urls []string) []FetchResult {

	ch := make(chan FetchResult, len(urls))
	var results []FetchResult

	for _, url := range urls {
		go fetch(url, ch)
	}

	for range urls {
		results = append(results, <-ch)
	}
	return results
}

func fetch(url string, ch chan FetchResult) {

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		ch <- FetchResult{URL: url, Error: fmt.Errorf("Couldn't fetch from URL %s: %v", url, err)}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- FetchResult{URL: url, Error: fmt.Errorf("Couldn't read the response from URL %s: %v", url, err)}
		return
	}

	ch <- FetchResult{
		URL:     url,
		Status:  resp.StatusCode,
		Content: string(body),
		Error:   nil,
	}
}
