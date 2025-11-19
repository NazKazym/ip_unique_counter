package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Population count (number of set bits in byte)
func popcount(b byte) int {
	b = (b & 0x55) + ((b >> 1) & 0x55)
	b = (b & 0x33) + ((b >> 2) & 0x33)
	b = (b & 0x0F) + ((b >> 4) & 0x0F)
	return int(b)
}

func startProgressReporter(linesProcessed, validIPs *atomic.Uint64) chan bool {
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		lastLines := uint64(0)
		lastTime := time.Now()

		for {
			select {
			case <-ticker.C:
				lines := linesProcessed.Load()
				valid := validIPs.Load()
				now := time.Now()

				linesDiff := lines - lastLines
				timeDiff := now.Sub(lastTime).Seconds()
				linesPerSec := float64(linesDiff) / timeDiff

				fmt.Printf("Progress: %d lines (%.1f M/s), %d valid IPs\n",
					lines, linesPerSec/1_000_000, valid)

				lastLines = lines
				lastTime = now
			case <-done:
				return
			}
		}
	}()
	return done
}
