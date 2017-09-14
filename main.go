package main

import "fmt"

func main() {
	events := make(chan int)

	go func() { events <- 2 }()

	for {
		responceTime := <-events
		logResponce(responceTime)
	}
}

func generateEvent()

func logResponce(time int) {
	fmt.Printf("Response took %d seconds\n", time)
}
