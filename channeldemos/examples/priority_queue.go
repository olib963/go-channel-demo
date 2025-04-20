package main

import (
	"fmt"
	"time"
)

func PriorityExample() {
	p, n, o := PriorityQueue[int]()

	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Printf("Sending normal message: %d\n", i)
			n <- i
			time.Sleep(333 * time.Millisecond)
		}
		println("Closing normal channel")
		close(n)
	}()

	go func() {
		time.Sleep(1 * time.Second)
		for i := 11; i <= 15; i++ {
			fmt.Printf("Sending priorty message: %d\n", i)
			p <- i
		}
		time.Sleep(1 * time.Second)
		for i := 16; i <= 20; i++ {
			fmt.Printf("Sending priorty message: %d\n", i)
			p <- i
		}
		println("Closing priority channel")
		close(p)
	}()

	results := make([]int, 0)
	for i := range o {
		results = append(results, i)
		time.Sleep(1 * time.Second)
	}

	fmt.Println(results)
}

func PriorityQueue[A any]() (chan<- A, chan<- A, <-chan A) {
	priorityQueue := make(chan A)
	normalQueue := make(chan A)
	outputQueue := make(chan A)
	go func() {
	Loop:
		for {
			select {
			// Try to get a priority message first
			case a, open := <-priorityQueue:
				if !open {
					break Loop
				}
				outputQueue <- a
				// If there isn't one, then race both channels. Then start the loop again
			default:
				select {
				case a, open := <-priorityQueue:
					if !open {
						break Loop
					}
					outputQueue <- a
				case a, open := <-normalQueue:
					if !open {
						break Loop
					}
					outputQueue <- a
				}
			}
		}
		// If either channel is closed, break out of the priority loop and simply drain the channels.
		for a := range priorityQueue {
			outputQueue <- a
		}
		for a := range normalQueue {
			outputQueue <- a
		}
		close(outputQueue)
	}()

	return priorityQueue, normalQueue, outputQueue
}
