package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
A demo of using select to read from multiple channels.
*/
func main() {
	greetings := make(chan string, 1)
	numbers := make(chan int, 1)

	wg := sync.WaitGroup{}
	wg.Add(2)

	println("Racing two channels...")

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		defer wg.Done()
		n := time.Duration(random.Intn(10))
		fmt.Printf("Thread 1 sleeping for %d seconds.\n", n)
		time.Sleep(n * time.Second)
		greetings <- "Hello"
	}()

	go func() {
		defer wg.Done()
		n := time.Duration(random.Intn(10))
		fmt.Printf("Thread 2 sleeping for %d seconds.\n", n)
		time.Sleep(n * time.Second)
		numbers <- 1
	}()

	race(greetings, numbers)

	wg.Wait()

	println("Racing a context with a channel...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// We never write to this channel
	willNeverWin := make(chan int)
	race(ctx.Done(), willNeverWin)

}

// Side note: we can use the `<-chan A` type to indicate that a function only reads from a channel.
// similarly, we can use the `chan<- A` type to indicate that a function only writes to a channel.
// writing to a read-only channel or reading from a write-only channel will cause a compile error.
func race[A, B any](c1 <-chan A, c2 <-chan B) {
	select {
	case a := <-c1:
		fmt.Printf("Channel 1 wins! %v\n", a)
	case b := <-c2:
		fmt.Printf("Channel 2 wins! %v\n", b)
	}
}
