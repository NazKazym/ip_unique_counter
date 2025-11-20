package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math/bits"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func (b *Bitmap) count() uint64 {
	numWorkers := runtime.NumCPU()
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

			var local uint64
			for j := start; j < end; j++ {
				local += uint64(bits.OnesCount64(b.data[j]))
			}
			results[workerID] = local
		}(i)
	}

	wg.Wait()

	var total uint64
	for _, c := range results {
		total += c
	}
	return total
}

func countWithBitmap(cfg Config) uint64 {
	bitmap := newBitmap()

	numWorkers := runtime.NumCPU() * 2
	if numWorkers > 128 {
		numWorkers = 128
	}

	fmt.Printf("Using %d parallel workers with single 512 MiB bitmap\n", numWorkers)

	file, err := os.Open(cfg.SourceURI)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	errorFile, err := os.Create("errors.log")
	if err != nil {
		log.Fatalf("cannot create errors.log: %v", err)
	}
	defer errorFile.Close()
	errorWriter := bufio.NewWriter(errorFile)
	defer errorWriter.Flush()

	var linesProcessed atomic.Uint64
	var validIPs atomic.Uint64

	start := time.Now()
	done := startProgressReporter(&linesProcessed, &validIPs, start)

	ipChans := make([]chan uint32, numWorkers)
	var wgWorkers sync.WaitGroup
	wgWorkers.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		ipChans[i] = make(chan uint32, 262_144)

		go func(ch <-chan uint32) {
			defer wgWorkers.Done()
			for ip := range ch {
				bitmap.set(ip)
			}
		}(ipChans[i])
	}

	bufferSizeBytes := cfg.Counter.BufferSizeMB * 1024 * 1024
	reader := bufio.NewReaderSize(file, bufferSizeBytes)

	const flushEvery = uint64(20_000_000)
	var localLines, localValid uint64
	var lineNumber uint64
	var wrongLines uint64

	for {
		lineBytes, err := reader.ReadSlice('\n')

		if len(lineBytes) > 0 {
			lineNumber++

			ip, err := parseIPv4Line(lineBytes)
			if err != nil {
				wrongLines++
				trimmed := trimRightSpaceCRLF(lineBytes)
				fmt.Fprintf(errorWriter, "%d | %q | %v\n", lineNumber, trimmed, err)
			} else {

				bitIdx := uint64(ip)
				wordIdx := bitIdx >> 6
				shardIdx := int(wordIdx % uint64(numWorkers))

				ipChans[shardIdx] <- ip
				localValid++
			}

			localLines++
			if localLines >= flushEvery {
				linesProcessed.Add(localLines)
				validIPs.Add(localValid)
				localLines, localValid = 0, 0
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil && !errors.Is(err, bufio.ErrBufferFull) {
			log.Fatal(err)
		}
	}

	if localLines > 0 {
		linesProcessed.Add(localLines)
		validIPs.Add(localValid)
	}

	for _, ch := range ipChans {
		close(ch)
	}
	wgWorkers.Wait()
	close(done)

	if wrongLines > 0 {
		fmt.Printf("\nWarning: %d invalid lines logged to errors.log\n", wrongLines)
	}

	fmt.Println("\nCounting unique IPs...")
	countStart := time.Now()
	count := bitmap.count()
	fmt.Printf("Count time: %s\n", time.Since(countStart))

	return count
}
