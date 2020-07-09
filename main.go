package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	flPort       int
	flPercent500 int

	latencyTreshold            int
	percentOverLatencyTreshold int
)

func init() {
	flag.IntVar(&flPort, "http-addr", 8080, "listen on http portrun on request (e.g. 8000)")
	flag.IntVar(&flPercent500, "percent500", -1, "percentage of requests that should get 500 error (0 <= n <= 100)")
	flag.Parse()

	// If no flag, try to get value from environment variable.
	if flPercent500 == -1 {
		envValue := os.Getenv("PERCENT_500_RESPONSES")
		if envValue == "" {
			flPercent500 = 0
		} else {
			var err error
			flPercent500, err = strconv.Atoi(envValue)
			if err != nil {
				log.Fatal("wrong value for percentage of 500 response")
			}
		}
	}
	if flPercent500 > 100 || flPercent500 < 0 {
		log.Fatalf("the percentage of 500 responses is not valid (0 <= n <= 100)")
	}

	// Latency tresholds.
	var err error
	value := os.Getenv("LATENCY_THRESHOLD")
	if value != "" {
		latencyTreshold, err = strconv.Atoi(value)
		if err != nil || latencyTreshold < 0 {
			log.Fatalf("wrong value for the latency treshold")
		}
	}

	value = os.Getenv("PERCENT_OVER_LATENCY_THRESHOLD")
	if value != "" {
		percentOverLatencyTreshold, err = strconv.Atoi(value)
		if err != nil || percentOverLatencyTreshold < 0 || percentOverLatencyTreshold > 100 {
			log.Fatalf("wrong value for the latency treshold")
		}
	}

}

func main() {
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flPort), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	s1 := rand.NewSource(time.Now().UTC().UnixNano())
	r1 := rand.New(s1)

	// Value 1 <= n <= 100
	randValue := 1 + r1.Intn(100)
	if randValue <= percentOverLatencyTreshold {
		duration := time.Duration(latencyTreshold + 10)
		time.Sleep(duration * time.Millisecond)
	}

	// Value 1 <= n <= 100
	randValue = 1 + rand.Intn(100)
	if randValue <= flPercent500 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	fmt.Fprintf(w, "successful request")
}
