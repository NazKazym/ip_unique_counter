package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func (b *Bitmap) count() uint64 {
	numWorkers := runtime.NumCPU() * 2
	chunkSize := len(b.data) / numWorkers

	var wg sync.WaitGroup
	results := make([]uint64, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			start := workerID * chunkSize
			end := start + chunkSize
			if workerID == numWorkers-1 {
				end = len(b.data)
			}

			var localCount uint64
			for j := start; j < end; j++ {
				localCount += uint64(popcount(b.data[j]))
			}
			results[workerID] = localCount
		}(i)
	}

	wg.Wait()

	var total uint64
	for _, c := range results {
		total += c
	}

	return total
}

type job struct {
	lines []string
}

// Strategy 2: Bitmap-based (for large files)
func countWithBitmap(cfg Config) uint64 {
	bitmap := newBitmap()
	numWorkers := runtime.NumCPU()

	fmt.Printf("Using %d workers with sharded bitmap (512 MB)\n", numWorkers)

	file, err := os.Open(cfg.SourceURI)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	jobs := make(chan job, numWorkers*2)
	var wg sync.WaitGroup
	var linesProcessed atomic.Uint64
	var validIPs atomic.Uint64

	// Calculate shard boundaries
	maxIP := uint64(1) << 32
	shardSize := maxIP / uint64(numWorkers)

	// Start workers - each handles specific IP range
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			shardStart := uint32(uint64(workerID) * shardSize)
			shardEnd := uint32(uint64(workerID+1) * shardSize)
			if workerID == numWorkers-1 {
				shardEnd = 0xFFFFFFFF
			}

			for j := range jobs {
				for _, line := range j.lines {
					ip, ok := parseIPv4Fast(line)
					if !ok {
						continue
					}

					// Only process IPs in this worker's range
					if ip >= shardStart && ip <= shardEnd {
						bitmap.set(ip)
						validIPs.Add(1)
					}
				}
				linesProcessed.Add(uint64(len(j.lines)))
			}
		}(i)
	}

	done := startProgressReporter(&linesProcessed, &validIPs)

	// Read file
	scanner := bufio.NewScanner(file)
	bufferSize := cfg.Counter.BufferSize
	batchSize := cfg.Counter.BatchSize
	scanner.Buffer(make([]byte, bufferSize), bufferSize*4)

	batch := make([]string, 0, batchSize)
	for scanner.Scan() {
		batch = append(batch, scanner.Text())

		if len(batch) >= batchSize {
			batchCopy := make([]string, len(batch))
			copy(batchCopy, batch)
			jobs <- job{lines: batchCopy}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		jobs <- job{lines: batch}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	close(jobs)
	wg.Wait()
	close(done)

	// Count unique IPs
	fmt.Println("\nCounting unique IPs...")
	countStart := time.Now()
	count := bitmap.count()
	fmt.Printf("Count time: %s\n", time.Since(countStart))

	return count
}
