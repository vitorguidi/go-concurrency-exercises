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
	"sync"
	"time"
)

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func Crawl(url string, depth int, wg *sync.WaitGroup, allowCrawlChannel *chan interface{}) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		// Do not remove the `go` keyword, as Crawl() must be
		// called concurrently
		go func() {
			<-*allowCrawlChannel
			fmt.Println(fmt.Sprintf("Crawling url %s at timestamp %s", url, time.Now()))
			Crawl(u, depth-1, wg, allowCrawlChannel)
		}()
	}
	return
}

func main() {
	var wg sync.WaitGroup
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

	wg.Add(1)
	Crawl("http://golang.org/", 4, &wg, &allowCrawlChannel)
	wg.Wait()
	cancel()
}
