/*
	Run length encoding for strings and []byte
*/
package rle

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type RLE struct {
	Flag       byte // defaults to 27 ESC
	CountBytes int  // 1 to 4 bytes defaults to 1
}

// countBytes return a valid CountByte value between 1 & 4 defaulting to 1 if not set
func (rle *RLE) countBytes() int {
	if rle == nil || rle.CountBytes < 1 {
		return 1
	} else if rle.CountBytes > 4 {
		return 4
	}
	return rle.CountBytes
}

// flag returns the flag byte, defaulting to 27 if not set
func (rle *RLE) flag() byte {
	if rle == nil {
		return 27
	}
	return rle.Flag
}

func (rle *RLE) EncodeString(in string) []byte {
	return rle.Encode([]byte(in))
}

func (rle *RLE) Encode(in []byte) []byte {
	size := len(in)
	if size == 0 {
		return []byte{}
	}

	countBytes := rle.countBytes()
	flag := rle.flag()
	var b bytes.Buffer
	cur := in[0]
	run := 0

	write := func() {
		switch {
		case run == 1:
			b.WriteByte(cur)

		case run <= countBytes+2: // Not worth encoding
			b.Write(bytes.Repeat([]byte{cur}, run))

		default: // Insert RLE block
			b.WriteByte(flag)
			count := make([]byte, 4)
			binary.BigEndian.PutUint32(count, uint32(run))
			b.Write(count[4-countBytes:])
			b.WriteByte(cur)
		}
	}

	for i := 0; i < size; i++ {
		c := in[i]
		if cur != c {
			write()
			cur = c
			run = 0
		}
		run++
	}
	write()

	return b.Bytes()
}

func (rle *RLE) DecodeString(in []byte) (string, error) {
	b, err := rle.Decode(in)
	return string(b), err
}

func (rle *RLE) Decode(b []byte) ([]byte, error) {
	countBytes := rle.countBytes()
	flag := rle.flag()

	var buf bytes.Buffer

	for len(b) > 0 {
		p := bytes.IndexByte(b, flag)
		if p == -1 {
			buf.Write(b)
			break
		}

		// Append the chars upto the flag to the output
		buf.Write(b[:p])

		// Skip the appended chars
		b = b[p+1:]

		// Check there are enough bytes left
		if len(b) < countBytes+1 {
			return b, errors.New("bad encoding")
		}

		// Get the count normalised as the max uint32
		count := binary.BigEndian.Uint32(append(
			bytes.Repeat([]byte{0x00}, 4-countBytes),
			b[:countBytes]...,
		))

		// Skip count
		b = b[countBytes:]

		// Get and skip char
		char := b[0]
		b = b[1:]

		// Append the repeated char to the output
		buf.Write(bytes.Repeat([]byte{char}, int(count)))
	}
	return buf.Bytes(), nil
}
