//////////////////////////////////////////////////////////////////////
//
// Your task is to change the code to limit the crawler to at most one
// page per second, while maintaining concurrency (in other words,
// Crawl() must be called concurrently)
//
// @hint: you can achieve this by adding 3 lines
//

package main

import (
	"context"
	"fmt"
	"time"
)

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func Crawl(url string, depth int, allowCrawlChannel *chan interface{}) {

	if depth <= 0 {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	for _, u := range urls {
		// Do not remove the `go` keyword, as Crawl() must be
		// called concurrently
		go func() {
			<-*allowCrawlChannel
			fmt.Println(fmt.Sprintf("Crawling url %s at timestamp %s", url, time.Now()))
			Crawl(u, depth-1, allowCrawlChannel)
		}()
	}
	return
}

func main() {
	allowCrawlChannel := make(chan interface{})
	ticker := time.NewTicker(1 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case <-ticker.C:
				allowCrawlChannel <- nil
			}
		}
	}()

	Crawl("http://golang.org/", 4, &allowCrawlChannel)
	cancel()
}
