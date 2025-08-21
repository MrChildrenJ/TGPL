package main

import (
	"fmt"
	"sync"
)

var (
	mutex   sync.Mutex
	wg      sync.WaitGroup
	visited = make(map[string]bool)
)

type Fetcher interface {
	// Fetch returns the body of URL and a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}

	mutex.Lock()
	if visited[url] {
		mutex.Unlock()
		return
	}
	visited[url] = true
	mutex.Unlock()

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	done := make(chan bool)
	count := 0 // numbers of go routines

	// 第一個for迴圈只是把工作分發給go routine，不代表工作做完
	// 會計算發了多少個工作出去 (count)
	for _, u := range urls {
		count++                   // Instant execution
		go func(nextUrl string) { // Instantly activate go routine, but...
			Crawl(nextUrl, depth-1, fetcher) // this will take a lot of time!!!!
			done <- true                     // wait for Crawl to finish
		}(u)
	}

	// Asynchronous
	// 等著接收 count 個 true
	for i := 0; i < count; i++ {
		<-done // Start waiting (for receiving) IMMEDIATELY !
	}
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
