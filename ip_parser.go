package main

import (
	"errors"
	"fmt"
)

var (
	ErrTooShort      = errors.New("too short for IPv4")
	ErrTooLong       = errors.New("too long for IPv4")
	ErrOctetTooLarge = errors.New("octet > 255")
	ErrTooManyDots   = errors.New("too many dots")
	ErrNotEnoughDots = errors.New("not exactly 4 octets")
)

func trimRightSpaceCRLF(b []byte) []byte {
	i := len(b)
	for i > 0 {
		c := b[i-1]
		if c == '\n' || c == '\r' || c == ' ' || c == '\t' {
			i--
		} else {
			break
		}
	}
	return b[:i]
}

func parseIPv4Line(b []byte) (uint32, error) {
	b = trimRightSpaceCRLF(b)

	ln := len(b)
	if ln > 15 {
		return 0, ErrTooLong
	}

	if ln < 7 {
		return 0, ErrTooShort
	}

	var ip uint32
	var octet uint32
	dots := 0

	for _, c := range b {
		switch {
		case c >= '0' && c <= '9':
			octet = octet*10 + uint32(c-'0')
			if octet > 255 {
				return 0, ErrOctetTooLarge
			}
		case c == '.':
			if dots >= 3 {
				return 0, ErrTooManyDots
			}
			ip = (ip << 8) | octet
			octet = 0
			dots++
		default:
			return 0, fmt.Errorf("invalid char %q", c)
		}
	}

	if dots != 3 {
		return 0, ErrNotEnoughDots
	}

	ip = (ip << 8) | octet
	return ip, nil
}
