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
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metricService *asyncinflux.AsyncClient
var promCounter prometheus.Counter
var promResponce prometheus.Summary

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file but ok - skip it")
	}

	events := make(chan int)
	go func() { generateEvents(events) }()

	initClients()

	for {
		responceTime := <-events
		logResponce(responceTime)
	}
}

func generateEvents(events chan<- int) {
	for {
		events <- int(rand.Int31n(10))
		time.Sleep(time.Second * 1)
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
	prometheus.MustRegister(promCounter)

	promResponce = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "test_responce_time",
		Help: "Time of test responce.",
	})
	prometheus.MustRegister(promResponce)

	http.Handle("/metrics", promhttp.Handler())
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) }()
}

func logResponce(time int) {
	fmt.Printf("Response took %d seconds\n", time)

	metricService.Send(asyncinflux.NewMetricDatum("test_measurement", map[string]string{}, map[string]interface{}{"response": time}))
	metricService.Send(asyncinflux.NewMetricDatum("test_measurement", map[string]string{}, map[string]interface{}{"count": 1}))
	promCounter.Inc()
	promResponce.Observe(float64(time))
}
