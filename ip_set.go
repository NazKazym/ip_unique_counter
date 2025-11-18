package main

type IPSet interface {
	AddUint32(ip uint32) bool
}

type mapIPSet struct {
	data map[uint32]struct{}
}

func newMapIPSet(capHint int) *mapIPSet {
	if capHint <= 0 {
		capHint = 1024
	}
	return &mapIPSet{
		data: make(map[uint32]struct{}, capHint),
	}
}

func (s *mapIPSet) AddUint32(ip uint32) bool {
	if _, exists := s.data[ip]; exists {
		return false
	}
	s.data[ip] = struct{}{}
	return true
}

type bitmapIPSet struct {
	bits []byte
}

func newBitmapIPSet() *bitmapIPSet {
	return &bitmapIPSet{
		bits: make([]byte, 1), // динамическое расширение
	}
}

func (s *bitmapIPSet) AddUint32(ip uint32) bool {
	byteIndex := ip / 8 // номер байта
	bitIndex := ip % 8  // номер бита в этом байте
	mask := byte(1 << bitIndex)

	// расширение массива при необходимости
	if int(byteIndex) >= len(s.bits) {
		newBits := make([]byte, byteIndex+1)
		copy(newBits, s.bits)
		s.bits = newBits
	}

	// если бит уже стоял
	if s.bits[byteIndex]&mask != 0 {
		return false
	}

	// ставим бит
	s.bits[byteIndex] |= mask
	return true
}
