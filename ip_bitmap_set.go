package main

type Bitmap struct {
	data []uint64
}

func newBitmap() *Bitmap {
	return &Bitmap{data: make([]uint64, 1<<26)}
}

func (b *Bitmap) set(ip uint32) {
	bitIdx := uint64(ip)
	wordIdx := bitIdx >> 6
	mask := uint64(1) << (bitIdx & 63)
	b.data[wordIdx] |= mask
}
