package main

import (
	"bufio"
	"context"
	"net"
	"os"
)

type UniqueIPCounter struct {
	cfg Config
}

func NewUniqueIPCounter(cfg Config) *UniqueIPCounter {
	return &UniqueIPCounter{cfg: cfg}
}

// CountUnique is still a simple in-memory unique counter (no buckets yet).
func (c *UniqueIPCounter) CountUnique(ctx context.Context) (int, error) {
	_ = ctx // not used yet

	f, err := os.Open(c.cfg.SourceURI)
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
