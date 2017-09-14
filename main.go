package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	asyncinflux "github.com/applift/async-influxdb-client"
	"github.com/joho/godotenv"
)

var metricService *asyncinflux.AsyncClient

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	events := make(chan int)
	go func() { generateEvents(events) }()

	initClients()

	for {
		responceTime := <-events
		logResponce(responceTime)
	}
}

func initClients() {
	// Setup defaults
	var err error
	metricService, err = asyncinflux.DefaultClient()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func generateEvents(events chan<- int) {
	for {
		events <- int(rand.Int31n(10))
		time.Sleep(time.Second * 2)
	}
}

func logResponce(time int) {
	fmt.Printf("Response took %d seconds\n", time)

	metricService.Send(asyncinflux.NewMetricDatum("test_measurement", map[string]string{}, map[string]interface{}{"response": time}))
}
