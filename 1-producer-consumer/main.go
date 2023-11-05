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
	"time"
)

func producer(stream Stream, tweets *chan *Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(*tweets)
			break
		}
		*tweets <- tweet
	}
}

func consumer(jobs <-chan *Tweet, done chan<- interface{}) {
	for {
		select {
		case tweet := <-jobs:
			if tweet == nil {
				done <- nil
				break
			}
			handleTweet(tweet)
		}
	}
}

func handleTweet(t *Tweet) {
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
	stream := GetMockStream()

	// Producer
	go producer(stream, &tweetsChan)

	// Consumer
	go consumer(tweetsChan, producingDone)
	<-producingDone

	fmt.Printf("Process took %s\n", time.Since(start))
}
