package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

// send message to channel using for loop and close channel
func sendMessage(ch chan<- string) {
	for i := 0; i < 20; i++ {
		ch <- fmt.Sprintf("Hello, World! using channel %d", i)
	}

	close(ch)
}

// print message from channel using range
func receiveMessage(ch <-chan string) {
	for message := range ch {
		fmt.Println(message)
	}
}

func sendData(ch chan<- string) {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; true; i++ {
		ch <- fmt.Sprintf("Hello, World! using channel %d", i)
		time.Sleep(time.Duration(randomizer.Intn(100)) * time.Millisecond)
	}
}

// print message from channel using select and timeout if no activity in 500ms
func retrieveData(ch <-chan string) {
	for {
		select {
		case message := <-ch:
			fmt.Println(message)
		case <-time.After(500 * time.Millisecond):
			fmt.Println("Timeout. No new message received.")
			return
		}
	}
}

func main() {
	runtime.GOMAXPROCS(2)

	// goroutine using range
	var messages = make(chan string)
	go sendMessage(messages)
	receiveMessage(messages)

	// goroutine using select
	messages = make(chan string)
	go sendData(messages)
	retrieveData(messages)
}
