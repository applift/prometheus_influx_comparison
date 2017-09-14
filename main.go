package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	asyncinflux "github.com/applift/async-influxdb-client"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

var metricService *asyncinflux.AsyncClient
var promCounter prometheus.Counter
var promResponce prometheus.Counter

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

	promCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_responce_counter",
		Help: "Number of test_responces.",
	})
	err = prometheus.Register(promCounter)
	if err != nil {
		log.Fatal(err.Error())
	}

	promResponce = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "test_responce_time",
		Help: "Time of test responce.",
	})
	err = prometheus.Register(promResponce)
	if err != nil {
		log.Fatal(err.Error())
	}

	http.Handle("/metrics", prometheus.Handler())
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) }()
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
	promCounter.Inc()
	promResponce.Add(float64(time))
}
