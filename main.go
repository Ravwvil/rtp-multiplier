package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

const (
	MinMultiplier = 1.0
	MaxMultiplier = 10000.0
	Port          = "64333"
)

type resp struct {
	Result float64 `json:"result"`
}

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func generateMultiplier(u, c float64) float64 {
	if c <= 0 {
		return MinMultiplier
	}
	if c > 1 {
		c = 1
	}
	massAtOne := 1.0 - c
	if massAtOne < 0 {
		massAtOne = 0
	}
	contUpper := 1.0 - c/MaxMultiplier // Верхняя граница непрерывной части распределения (CDF)
	if u < massAtOne {                 // Точечная масса в M=1
		return MinMultiplier
	}
	if u < contUpper { // Непрерывная часть распределения
		m := c / (1.0 - u)
		return clamp(m, MinMultiplier, MaxMultiplier)
	}
	return MaxMultiplier // Точечная масса в M=10000
}

func main() {
	rtpFlag := flag.Float64("rtp", -1, "target RTP (0 < rtp <= 1.0)")
	flag.Parse()
	if *rtpFlag <= 0 || *rtpFlag > 1.0 {
		fmt.Fprintln(os.Stderr, "error: -rtp must be provided and 0 < rtp <= 1.0")
		os.Exit(2)
	}
	rtp := *rtpFlag

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		u := rand.Float64()
		val := generateMultiplier(u, rtp)
		json.NewEncoder(w).Encode(resp{Result: val})
	})

	log.Printf("Multiplier service listening :%s — rtp=%.6f", Port, rtp)
	log.Fatal(http.ListenAndServe(":"+Port, nil))
}
