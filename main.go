package main

import (
	"context"
	"fmt"
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

	counter := NewUniqueIPCounter(cfg)

	unique, err := counter.CountUnique(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Unique IPs:", unique)
}
