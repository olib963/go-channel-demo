package main

import "time"

func main() {
	println("Running Priority Queue example...")
	PriorityExample()

	time.Sleep(10 * time.Second)
	println("\n\nRunning Rate limiter example...")
	LimiterExample()
}
