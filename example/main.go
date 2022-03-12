package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/kgolding/go-rle"
)

func main() {
	// Our chosen flag
	flag := byte(27)

	// Generate some random but repeated data
	rand.Seed(time.Now().Unix())
	b := make([]byte, 0)
	r := make([]byte, 2)
	for i := 0; i < 100; i++ {
	getRandData:
		rand.Read(r)
		if r[0] == flag {
			// Our data must not include the flag byte
			goto getRandData
		}
		b = append(b, bytes.Repeat([]byte{r[0]}, int(r[1]))...)
	}

	fmt.Printf("Raw data size     % 8d\n", len(b))

	// Create an RLE encoder/decoder with the default settings
	enc := &rle.RLE{}

	// Encode the data to out
	out := enc.Encode(b)

	fmt.Printf("RLE encoded size  % 8d\n", len(out))

	// Decode the encoded data
	in, err := enc.Decode(out)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	fmt.Printf("RLE decoded size  % 8d\n", len(in))
}
