package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
)

const (
	MinMultiplier = 1.0
	MaxMultiplier = 10000.0
	L             = MaxMultiplier - MinMultiplier
	Port          = "64333"
)

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func computeThreshold(rtp float64) float64 {
	// m0 = √(2 * rtp * L + 1)
	val := math.Sqrt(2.0*rtp*L + MinMultiplier*MinMultiplier)
	return clamp(val, MinMultiplier, MaxMultiplier)
}

type resp struct {
	Result float64 `json:"result"`
}

func main() {
	rtpFlag := flag.Float64("rtp", -1, "target RTP (0 < rtp <= 1.0). Example: -rtp=0.96")
	flag.Parse()

	fmt.Printf("Parsed rtp value: %f\n", *rtpFlag)
	rtp := *rtpFlag
	if rtp <= 0 || rtp > 1.0 {
		fmt.Fprintln(os.Stderr, "error: -rtp must be provided and 0 < rtp <= 1.0 (e.g. -rtp=0.96)")
		os.Exit(2)
	}

	m0 := computeThreshold(rtp)
	if m0 >= MaxMultiplier {
		log.Printf("warning: computed m0 >= %.0f; clamped to %.0f — requested rtp may be unreachable under uniform-x model", MaxMultiplier, MaxMultiplier)
	}

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := resp{Result: m0}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("encode error: %v", err)
		}
	})

	addr := ":" + Port
	log.Printf("Multiplier service listening %s — returning constant m0=%.6f for rtp=%.6f", addr, m0, rtp)
	log.Fatal(http.ListenAndServe(addr, nil))
}
