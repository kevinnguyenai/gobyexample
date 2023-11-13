// We can use channels to synchronize the execution state.
// There is an example here. The method of blocking reception is used to wait for another coroutine to complete.
// If you need to wait for multiple coroutines, [Waitgroup] (Waitgroups) is a better choice.

package main

import (
	"fmt"
	"time"
)

// We will run this function in the coroutine.
// `Done` channels will be used to notify other corporate functions that have been completed.
func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	// Send a value to inform us that we are completed.
	done <- true
}

func main() {

	// Run a worker coroutine and give a channel for notification.
	done := make(chan bool, 1)
	go worker(done)

	// The program will always be blocked until the notice sent by the channel uses channel.
	<-done
}
