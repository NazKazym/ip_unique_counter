package main

// Parses IP directly from string to uint32
func parseIPv4Fast(s string) (uint32, bool) {
	var a, b, c, d int
	var n int
	parts := 0

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch >= '0' && ch <= '9' {
			n = n*10 + int(ch-'0')
			if n > 255 {
				return 0, false
			}
		} else if ch == '.' {
			switch parts {
			case 0:
				a = n
			case 1:
				b = n
			case 2:
				c = n
			default:
				return 0, false
			}
			parts++
			n = 0
		} else if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			// Allow whitespace after IP
			break
		} else {
			// Invalid character
			return 0, false
		}
	}

	if parts != 3 {
		return 0, false
	}
	d = n

	// Single check for all octets
	if a|b|c|d > 255 {
		return 0, false
	}

	ip := uint32(a)<<24 | uint32(b)<<16 | uint32(c)<<8 | uint32(d)
	return ip, true
}
