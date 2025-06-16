package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main11() {
	start := time.Now()
	ch := make(chan string)
	
	urls := os.Args[1:]
	if len(urls) == 0 {
		urls = []string{
			"https://google.com",
			"https://youtube.com", 
			"https://facebook.com",
			"https://baidu.com",
			"https://wikipedia.org",
			"https://qq.com",
			"https://taobao.com",
			"https://amazon.com",
			"https://twitter.com",
			"https://instagram.com",
			"https://sohu.com",
			"https://yahoo.com",
			"https://jd.com",
			"https://reddit.com",
			"https://weibo.com",
			"https://360.cn",
			"https://sina.com.cn",
			"https://live.com",
			"https://netflix.com",
			"https://microsoft.com",
			
			// Add some potentially problematic URLs
			"https://nonexistent-website-12345.com",
			"https://httpstat.us/500", // Returns 500 error
			"https://httpbin.org/delay/10", // 10 second delay
			"https://httpstat.us/timeout", // Simulates timeout
		}
		fmt.Printf("Testing %d popular websites...\n", len(urls))
	}
	
	for _, url := range urls {
		go fetch(url, ch) // start a goroutine
	}
	
	for range urls {
		fmt.Println(<-ch) // receive from channel ch
	}
	
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	
	// resp, err := http.Get(url)  // Use default client, no time out

	// Create HTTP client with timeout
	client := &http.Client{			// client is a pointer!!
		Timeout: 10 * time.Second, 	// 10 second timeout
		// https://pkg.go.dev/net/http#Client
	}
	
	resp, err := client.Get(url)
	if err != nil {
		secs := time.Since(start).Seconds()
		ch <- fmt.Sprintf("%.2fs  ERROR   %s: %v", secs, url, err)
		return
	}
	
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	
	if err != nil {
		secs := time.Since(start).Seconds()
		ch <- fmt.Sprintf("%.2fs  ERROR   while reading %s: %v", secs, url, err)
		return
	}
	
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

/*
Testing 24 popular websites...
0.05s  ERROR   https://nonexistent-website-12345.com: Get "https://nonexistent-website-12345.com": dial tcp: lookup nonexistent-website-12345.com: no such host
0.30s    17675  https://google.com
0.30s    91470  https://wikipedia.org
0.32s    82012  https://facebook.com
0.32s     2540  https://twitter.com
0.32s   238832  https://microsoft.com
0.34s       23  https://yahoo.com
0.34s       25  https://httpstat.us/500
0.35s        0  https://httpstat.us/timeout
0.54s   701391  https://live.com
0.56s     6591  https://amazon.com
0.78s   454758  https://instagram.com
0.79s   619144  https://youtube.com
0.98s   647737  https://reddit.com
1.00s  3045913  https://netflix.com
1.07s    81203  https://taobao.com
1.32s     2381  https://baidu.com
1.36s    93298  https://360.cn
1.40s        0  https://sohu.com
1.43s   395454  https://sina.com.cn
2.14s      326  https://qq.com
2.43s     8393  https://weibo.com
3.43s    17824  https://jd.com
10.44s      323  https://httpbin.org/delay/10
10.44s elapsed
*/