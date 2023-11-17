package main

import (
	"fmt"
	"sync"
	"time"
)

// TODO for loop demo
// TODO broadcast demo
// TODO select demo
// TODO actors
/*
A small demo of buffered channels, these are channels that can hold a certain number of messages before blocking.
*/
func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	ints := make(chan int, 5)

	go func() {
		defer wg.Done()
		for i := 1; i <= 15; i++ {
			ints <- i
			fmt.Printf("Sent %d to channel.\n", i)
		}
		println("Closing Channel...") // TODO should the closing be a separate demo
		close(ints)                   // We can close a channel to indicate that we are done sending messages.
	}()

	go func() {
		defer wg.Done()
		total := 0
		for i := range ints { // We can use range to iterate over a channel until it is closed.
			fmt.Printf("Received %d from channel.\n", i)
			time.Sleep(time.Second)
			total += i
		}
		fmt.Printf("Total: %d\n", total)
	}()

	wg.Wait()
}
