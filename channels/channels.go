package main

import (
	"fmt"
)

const channelBufSize = 4
const numDataElements = 100

func main() {
	muxChannels := make([]chan int, 4)
	demuxChannels := make([]chan int, 4)

	// Create channels
	drainChan := make(chan int, channelBufSize)
	for i := 0; i < 4; i++ {
		muxChannels[i] = make(chan int, channelBufSize)
		demuxChannels[i] = make(chan int, channelBufSize)
	}

	// Worker routines, each with an individual channel
	for i := 0; i < 4; i++ {
		go func(input <-chan int, output chan<- int) {
			// Get data, blocking while input channel open
			for i := range input {
				output <- i
			}
			close(output)
		}(muxChannels[i], demuxChannels[i])
	}

	// data source
	go func() {
		// Communicate with channels
		for i := 0; i < numDataElements; i++ {
			muxChannels[i%4] <- i
		}
		// Close mux channels
		for _, muxChan := range muxChannels {
			close(muxChan)
		}
	}()

	// demultiplex
	go func(multiplex []chan int, output chan<- int) {
		var ok1, ok2, ok3, ok4 bool
		for {
			var i1, i2, i3, i4 int
			select {
			case i1, ok1 = <-multiplex[0]:
				if ok1 {
					output <- i1
				}
			case i2, ok2 = <-multiplex[1]:
				if ok2 {
					output <- i2
				}
			case i3, ok3 = <-multiplex[2]:
				if ok3 {
					output <- i3
				}
			case i4, ok4 = <-multiplex[3]:
				if ok4 {
					output <- i4
				}
			}

			// When no channels are open, close output and exit
			if !ok1 && !ok2 && !ok3 && !ok4 {
				close(output)
				break
			}
		}
	}(demuxChannels, drainChan)

	// Final destination, data drain
	received := make([]int, 0, numDataElements)
	for i := range drainChan {
		received = append(received, i)
	}

	fmt.Println(received)
	fmt.Println(len(received))
}

// vim: ts=4 sw=4 ai
