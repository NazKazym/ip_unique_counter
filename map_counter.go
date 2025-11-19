package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

// Map-based counter for small files
type MapCounter struct {
	shards [256]*mapShard
}

type mapShard struct {
	mu  sync.Mutex
	ips map[uint32]struct{}
}

func newMapCounter() *MapCounter {
	counter := &MapCounter{}
	for i := 0; i < 256; i++ {
		counter.shards[i] = &mapShard{
			ips: make(map[uint32]struct{}),
		}
	}
	return counter
}

func (m *MapCounter) add(ip uint32) {
	shardIdx := ip % 256
	s := m.shards[shardIdx]
	s.mu.Lock()
	s.ips[ip] = struct{}{}
	s.mu.Unlock()
}

func (m *MapCounter) count() uint64 {
	var total uint64
	for i := 0; i < 256; i++ {
		m.shards[i].mu.Lock()
		total += uint64(len(m.shards[i].ips))
		m.shards[i].mu.Unlock()
	}
	return total
}

// Strategy 1: Map-based (for small files)
func countWithMap(cfg Config) uint64 {
	counter := newMapCounter()
	numWorkers := runtime.NumCPU()

	fmt.Printf("Using %d workers with sharded map\n", numWorkers)

	file, err := os.Open(cfg.SourceURI)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	jobs := make(chan job, numWorkers*2)
	var wg sync.WaitGroup
	var linesProcessed atomic.Uint64
	var validIPs atomic.Uint64

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := range jobs {
				for _, line := range j.lines {
					ip, ok := parseIPv4Fast(line)
					if !ok {
						continue
					}

					counter.add(ip)
					validIPs.Add(1)
				}
				linesProcessed.Add(uint64(len(j.lines)))
			}
		}()
	}

	done := startProgressReporter(&linesProcessed, &validIPs)

	// Read file
	bufferSize := cfg.Counter.BufferSize * 1024 * 1024
	batchSize := cfg.Counter.BatchSize
	scanner := bufio.NewScanner(file)
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

	return counter.count()
}
