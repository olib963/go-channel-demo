package main

import (
	"sync"
	"time"
)

/*
A demo of an unbuffered channel passing messages between two threads. It demonstrates the blocking nature of
channels on both the sender and receiver.
*/
func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	channel := make(chan string)

	go func() {
		defer wg.Done()

		println("Sending to channel...")
		// This will block until the receiver is ready to receive the message.
		channel <- "Hello World!"
		println("Sent to channel.")

		time.Sleep(25 * time.Second)

		println("Sending to channel again...")
		channel <- "Goodbye World!"
		println("Sent to channel again.")
	}()

	time.Sleep(10 * time.Second)

	println("Receiving from channel...")
	println("Received from channel: " + <-channel)

	println("Receiving from channel again...")
	// This will block until the sender is ready to send the message.
	println("Received from channel: " + <-channel)

	wg.Wait()
}
