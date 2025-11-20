package main

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

const progressStep = 1_000_000

func startProgressReporter(
	linesProcessed, validIPs *atomic.Uint64,
	startTime time.Time,
) chan<- struct{} {

	done := make(chan struct{}, 1)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		var lastLogged uint64

		for {
			select {
			case <-ticker.C:
				totalLines := linesProcessed.Load()
				if totalLines == 0 {
					continue
				}

				if totalLines < lastLogged+progressStep {
					continue
				}

				totalValid := validIPs.Load()
				now := time.Now()
				printProgress(startTime, now, totalLines, totalValid)

				lastLogged = totalLines - (totalLines % progressStep)

			case <-done:
				totalLines := linesProcessed.Load()
				totalValid := validIPs.Load()
				totalTime := time.Since(startTime)

				fmt.Printf("\n\nFinished!\n")
				fmt.Printf("   Total lines processed : %s (%.3f B)\n",
					formatUint64(totalLines), float64(totalLines)/1e9)
				fmt.Printf("   Valid IPv4 addresses  : %s (%.3f M)\n",
					formatUint64(totalValid), float64(totalValid)/1e6)
				fmt.Printf("   Total time            : %s\n",
					totalTime.Round(time.Millisecond))
				fmt.Printf("   Average speed         : %.2f M lines/sec\n",
					float64(totalLines)/totalTime.Seconds()/1e6)
				return
			}
		}
	}()

	return done
}

func printProgress(start, now time.Time, lines, valid uint64) {
	elapsed := now.Sub(start)
	if elapsed <= 0 {
		return
	}

	speed := float64(lines) / elapsed.Seconds()

	fmt.Printf(
		"\rProgress :: %s lines | %6.2f M/s | %s IPs | elapsed %s",
		formatUint64(lines),
		speed/1e6,
		formatUint64(valid),
		elapsed.Round(100*time.Millisecond),
	)
}

func formatUint64(n uint64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	var parts []string
	s := fmt.Sprintf("%d", n)
	for i := len(s); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{s[start:i]}, parts...)
	}
	return strings.Join(parts, ",")
}
