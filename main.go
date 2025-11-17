package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
)

func main() {
	configPath := defaultConfigFileName
	if len(os.Args) >= 2 {
		configPath = os.Args[1]
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	unique, err := countUniqueIPsInMemory(context.Background(), cfg.SourceURI)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Unique IPs:", unique)
}

func countUniqueIPsInMemory(_ context.Context, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	seen := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		ip := net.ParseIP(line)
		if ip == nil {
			continue
		}
		seen[line] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return len(seen), nil
}
