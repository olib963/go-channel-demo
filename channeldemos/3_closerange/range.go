package main

import (
	"fmt"
	"sync"
	"time"
)

/*
A small demo of closing a channel when you are done with it and using range to iterate over the messages in the channel.
*/
func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	ints := make(chan int, 5)

	go func() {
		defer wg.Done()

		subGroup := sync.WaitGroup{}
		subGroup.Add(5)

		// Send all the numbers between 1 and 15 to the channel in 5 goroutines.
		for threadIndex := 1; threadIndex <= 5; threadIndex++ {
			go func(i int) {
				defer subGroup.Done()
				for j := 1; j <= 3; j++ {
					data := 3*i - j + 1
					ints <- data
					fmt.Printf("Sent %d to channel.\n", data)
				}
			}(threadIndex)
		}

		subGroup.Wait()
		println("Closing Channel...")
		close(ints) // We can close a channel to indicate that we are done sending messages.
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
