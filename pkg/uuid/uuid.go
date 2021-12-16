package uuid

import (
	"encoding/hex"
)

// Nil is nil value of UUID.
var Nil UUID

// A UUID is a 128 bit (16 byte) Universal Unique Identifier
// as defined in RFC 4122.
type UUID [16]byte

// String returns uuid as a formatted string.
func (id UUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], id)
	return string(buf[:])
}

func encodeHex(dst []byte, id UUID) {
	hex.Encode(dst, id[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], id[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], id[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], id[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], id[10:])
}
