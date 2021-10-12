package main

import (
	"fmt"
	"time"
)

func producer(stream Stream) <-chan Tweet {
	ch := make(chan Tweet)

	go func() {
		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				close(ch)
				break
			}

			ch <- *tweet
		}
	}()

	return ch
}

func consumer(tweets <-chan Tweet) {
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
			continue
		}

		fmt.Println(t.Username, "\tdoes not tweet about golang")
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Modification starts from here
	// Hint: this can be resolved via channels
	// Producer
	tweets := producer(stream)
	// Consumer
	consumer(tweets)

	fmt.Printf("Process took %s\n", time.Since(start))
}
