package main

import (
	"net"
	"strings"
)

// IPValidator validates and optionally normalizes an IP line.
type IPValidator interface {
	Validate(line string) (string, bool)
}

// DefaultIPValidator parses IPs using net.ParseIP.
type DefaultIPValidator struct{}

// Validate returns normalized IP string if valid, otherwise false.
func (v *DefaultIPValidator) Validate(line string) (string, bool) {
	s := strings.TrimSpace(line)
	if s == "" {
		return "", false
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return "", false
	}
	// You could normalize to ip.String() if you want canonical forms.
	return s, true
}
