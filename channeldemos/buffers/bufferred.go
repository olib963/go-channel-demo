package main

import (
	"fmt"
	"sync"
	"time"
)

/*
A small demo of buffered channels, these are channels that can hold a certain number of messages before blocking.
*/
func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	ints := make(chan int, 5)

	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			ints <- i
			fmt.Printf("Sent %d to channel.\n", i)
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Second)
		for i := 1; i <= 10; i++ {
			fmt.Printf("Received %d.\n", <-ints)
		}
	}()

	wg.Wait()
}
