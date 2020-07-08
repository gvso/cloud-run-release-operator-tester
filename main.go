package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var (
	flPort       int
	flPercent500 int

	minSleep int
	maxSleep int
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
}

func main() {
	log.Println(flPercent500)
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flPort), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Value 1 <= n <= 100
	randValue := 1 + rand.Intn(100)
	if randValue <= flPercent500 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}

	fmt.Fprintf(w, "successful request")
}
