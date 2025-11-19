package main

const ( // Lines per batch
	bitmapSize = 512 * 1024 * 1024 // 512MB for all IPv4
)

// Bitmap for all IPv4 addresses
type Bitmap struct {
	data []byte
}

func newBitmap() *Bitmap {
	return &Bitmap{
		data: make([]byte, bitmapSize),
	}
}

// Set bit for IP (no locks needed with sharding)
func (b *Bitmap) set(ip uint32) {
	byteIdx := ip / 8
	bitIdx := ip % 8
	b.data[byteIdx] |= 1 << bitIdx
}
