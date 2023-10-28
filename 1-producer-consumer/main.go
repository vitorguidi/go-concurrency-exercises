//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, tweets *chan *Tweet, producingDone *chan interface{}, wg *sync.WaitGroup) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(*tweets)
			*producingDone <- nil
			break
		}
		wg.Add(1)
		*tweets <- tweet
	}
}

func consumer(t *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	if t.IsTalkingAboutGo() {
		fmt.Println(t.Username, "\ttweets about golang")
	} else {
		fmt.Println(t.Username, "\tdoes not tweet about golang")
	}
}

func main() {
	start := time.Now()
	tweetsChan := make(chan *Tweet)
	producingDone := make(chan interface{})
	wg := sync.WaitGroup{}
	stream := GetMockStream()

	// Producer
	go producer(stream, &tweetsChan, &producingDone, &wg)

	// Consumer
	for x := range tweetsChan {
		go consumer(x, &wg)
	}
	<-producingDone
	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
