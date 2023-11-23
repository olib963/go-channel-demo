package main

import (
	"fmt"
	"sync"
	"time"
)

/*
A demo of closing a channel when you are done with it and using range to iterate over the messages in the channel.
*/
func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	ints := make(chan int, 5)

	go func() {
		defer wg.Done()

		subGroup := sync.WaitGroup{}
		subGroup.Add(5)

		// Send all the numbers between 1 and 15 to the channel in 5 goroutines.
		for threadIndex := 1; threadIndex <= 5; threadIndex++ {
			go func(threadIndex int) {
				defer subGroup.Done()
				for n := 1; n <= 3; n++ {
					data := 3*threadIndex - n + 1
					ints <- data
					fmt.Printf("Sent %d to channel from thread %d.\n", data, threadIndex)
				}
			}(threadIndex)
		}

		subGroup.Wait()
		println("Closing Channel...")
		close(ints) // We can close a channel to indicate that we are done sending messages.
	}()

	total := 0
	for i := range ints { // We can use range to iterate over a channel until it is closed.
		fmt.Printf("Received %d from channel.\n", i)
		time.Sleep(time.Second)
		total += i
	}
	fmt.Printf("Total: %d\n", total)

	last, ok := <-ints // We can check if a channel is closed by reading from it, you'll get the zero false and false if it is closed.
	fmt.Printf("Done %d, %v.\n", last, ok)
	wg.Wait()
}
