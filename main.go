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

	latency99 int
	latency95 int
	latency50 int
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

	// Latency thresholds.
	var err error
	value := os.Getenv("LATENCY_P99")
	if value != "" {
		latency99, err = strconv.Atoi(value)
		if err != nil {
			log.Fatalf("wrong value for the latency 99th percentile")
		}
	}

	value = os.Getenv("LATENCY_P95")
	if value != "" {
		latency95, err = strconv.Atoi(value)
		if err != nil {
			log.Fatalf("wrong value for the latency 95th percentile")
		}
	}

	value = os.Getenv("LATENCY_P50")
	if value != "" {
		latency50, err = strconv.Atoi(value)
		if err != nil {
			log.Fatalf("wrong value for the latency 50th percentile")
		}
	}

	if latency50 > latency95 || latency95 > latency99 {
		log.Fatal("99th percentile >= 95 percentile >= 50th percentile")
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
	if randValue <= 1 {
		time.Sleep(time.Duration(latency99) * time.Millisecond)
	} else if randValue <= 5 {
		time.Sleep(time.Duration(latency95) * time.Millisecond)
	} else if randValue <= 50 {
		time.Sleep(time.Duration(latency50) * time.Millisecond)
	}

	// Value 1 <= n <= 100
	randValue = 1 + rand.Intn(100)
	if randValue <= flPercent500 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	// 200 and 201.
	response := 200 + rand.Intn(2)
	w.WriteHeader(response)

	fmt.Fprintf(w, "successful request")
}
