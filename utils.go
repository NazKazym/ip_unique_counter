package main

import "os"

// средний размер строки "A.B.C.D\n" ~ 12–18 байт.
// возьмём, например, 16.
const avgBytesPerLine = 16

func approximateItemsFromFile(path string) uint64 {
	fi, err := os.Stat(path)
	if err != nil {
		// если не смогли получить размер — вернём 0, дальше будем жить с map
		return 0
	}
	size := fi.Size()
	if size <= 0 {
		return 0
	}
	return uint64(size) / avgBytesPerLine
}
