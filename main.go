package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: unique_ip_count <path-to-ip-file>")
		os.Exit(1)
	}

	path := os.Args[1]
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		os.Exit(1)
	}
	defer f.Close()

	seen := make(map[string]struct{})

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		ip := net.ParseIP(line)
		if ip == nil {
			continue // skip invalid lines
		}
		seen[line] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Read error:", err)
		os.Exit(1)
	}

	fmt.Println("Unique IPs:", len(seen))
}
