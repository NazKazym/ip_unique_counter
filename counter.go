package main

import (
	"bufio"
	"context"
	"os"
)

type UniqueIPCounter struct {
	cfg Config
}

func NewUniqueIPCounter(cfg Config) *UniqueIPCounter {
	return &UniqueIPCounter{cfg: cfg}
}
func (c *UniqueIPCounter) CountUnique(ctx context.Context) (int, error) {
	// 1) Примерно оцениваем количество строк/айпишников по размеру файла.
	approxItems := approximateItemsFromFile(c.cfg.SourceURI)

	f, err := os.Open(c.cfg.SourceURI)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// 2) Выбираем реализацию (map/bitmap) уже с учётом approxItems.
	ipSet := newIPSetFromConfig(c.cfg.Counter, approxItems)

	scanner := bufio.NewScanner(f)
	unique := 0

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return unique, ctx.Err()
		default:
		}

		line := scanner.Text()

		ip, ok := ParseIPv4(line)
		if !ok {
			continue
		}

		if ipSet.AddUint32(ip) {
			unique++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return unique, nil
}
