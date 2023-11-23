package main

import (
	"fmt"
	"sync"
	"time"
)

/*
A demo of buffered channels, these are channels that can hold a certain number of messages before blocking.
*/
func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	ints := make(chan int, 5)

	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			// This will not block until the channel is full which will happen after 5 messages are sent.
			ints <- i
			fmt.Printf("Sent %d to channel.\n", i)
		}
	}()

	time.Sleep(10 * time.Second)
	for i := 1; i <= 10; i++ {
		fmt.Printf("Received %d.\n", <-ints)
	}

	wg.Wait()
}
