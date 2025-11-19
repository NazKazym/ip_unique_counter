package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	configPath := defaultConfigFileName

	cfg, err := LoadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}
	fmt.Printf("Config properties:")
	fmt.Printf("Buffer size: %d\nMB", cfg.Counter.BufferSize)
	fmt.Printf("Batch size: %d\n lines", cfg.Counter.BatchSize)
	start := time.Now()

	count := countUniqueIPs(cfg)

	elapsed := time.Since(start)
	fmt.Printf("\nUnique IPv4 addresses: %d\n", count)
	fmt.Printf("Time elapsed: %s\n", elapsed)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Memory used: %.2f MB\n", float64(m.Alloc)/1024/1024)
}

func countUniqueIPs(cfg Config) uint64 {
	// Get file size
	filePath := cfg.SourceURI
	bitmapThresholdBytes := cfg.BitmapThresholdMB
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileInfo.Size()

	fmt.Printf("File size: %.2f MB\n", float64(fileSize)/(1024*1024))

	// Choose strategy based on file size
	if fileSize < int64(bitmapThresholdBytes) {
		fmt.Printf("Using MAP strategy (file < %d MB)\n", bitmapThresholdBytes)
		return countWithMap(cfg)
	} else {
		fmt.Printf("Using BITMAP strategy (file >= %d MB)\n", bitmapThresholdBytes)
		return countWithBitmap(cfg)
	}
}
