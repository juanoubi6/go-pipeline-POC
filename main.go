package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var wgGorutines sync.WaitGroup

// Expected result: the counter value will increase in each goroutine.
func main() {
	println("Number of processors: " + strconv.Itoa(runtime.NumCPU()))
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create pipeline channels
	var startCh = make(chan int)
	var channel2 = make(chan int)
	var channel3 = make(chan int)
	var channel4 = make(chan int)
	var finalCh = make(chan int)

	// Create a WaitGroup to wait for all goroutine to finish before exiting
	wgGorutines.Add(4)

	go worker(startCh, channel2, 1)
	go worker(channel2, channel3, 2)
	go worker(channel3, channel4, 3)
	go worker(channel4, finalCh, 4)

	// Start the pipeline by sending a value to the first channel of the chain
	startCh <- 0

	// Block until a message is published in the final channel
	finalValue := <-finalCh
	println("Received value from last channel: ", finalValue)

	println("Number of gorutines executing: " + strconv.Itoa(runtime.NumGoroutine()))
}

func worker(receiveCh chan int, sendCh chan int, id int) {
	defer wgGorutines.Done()

	// Consume from the receive channel. If there is no element in it, block until someone publishes something.
	counterValue := <-receiveCh
	fmt.Printf("Worker with id %v is running \n", id)
	println("Got value from pipeline: ", counterValue)
	counterValue++
	println("Increasing value from pipeline: ", counterValue)

	// Send increased value to the send channel so the flow can continue.
	fmt.Printf("Worker with id %v closing... \n", id)
	sendCh <- counterValue

}
