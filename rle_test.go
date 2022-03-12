package rle

import (
	"bytes"
	"testing"
)

var tests = []struct {
	encoded []byte
	decoded []byte
}{
	{[]byte{}, []byte("")},
	{[]byte("abc"), []byte("abc")},
	{[]byte("QWERTYUIOP!\"£$%^&*()ASDFGHJKLZXCVBNM<>+}~?_P{:?>, \""), []byte("QWERTYUIOP!\"£$%^&*()ASDFGHJKLZXCVBNM<>+}~?_P{:?>, \"")},
	{[]byte{31}, []byte{31}},
	{[]byte{31, 31}, []byte{31, 31}},
	{[]byte{31, 31, 31}, []byte{31, 31, 31}},
	{[]byte{27, 4, 31}, []byte{31, 31, 31, 31}},
	{[]byte{27, 10, 'X'}, []byte("XXXXXXXXXX")},
	{[]byte{27, 10, 'X', 27, 4, 'M'}, []byte("XXXXXXXXXXMMMM")},
	{[]byte{'A', 'B', 27, 4, 'C'}, []byte("ABCCCC")},
	{[]byte{'A', 'B', 27, 4, 'C', 'D'}, []byte("ABCCCCD")},
}

func TestDecode(t *testing.T) {
	for _, flag := range []byte{0, 1, 2, 27, 255} {
		rle := RLE{
			Flag:       flag,
			CountBytes: 1,
		}
		for _, test := range tests {
			in, err := rle.Decode(bytes.ReplaceAll(test.encoded, []byte{27}, []byte{flag}))
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(in, test.decoded) != 0 {
				t.Errorf("expected '% X' got '% X'", test.decoded, in)
			}
		}
	}
}

func TestEncode(t *testing.T) {
	for _, flag := range []byte{0, 1, 2, 27, 255} {
		rle := RLE{
			Flag:       flag,
			CountBytes: 1,
		}
		for _, test := range tests {
			b := rle.Encode(test.decoded)
			if bytes.Compare(b, bytes.ReplaceAll(test.encoded, []byte{27}, []byte{flag})) != 0 {
				t.Errorf("expected:\n\t'% X' got\n\t'% X'", test.encoded, b)
			}
		}
	}
}

func TestCountLengths(t *testing.T) {
	data := bytes.Repeat([]byte{65}, 10)
	for _, cl := range []int{1, 2, 3, 4} {
		rle := &RLE{
			Flag:       27,
			CountBytes: cl,
		}
		expectedLen := 2 + cl
		b := rle.Encode(data)
		if len(b) != expectedLen {
			t.Errorf("expected encoded length to be %d got %d: %X", expectedLen, len(b), b)
		}
		x, err := rle.Decode(b)
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(data, x) != 0 {
			t.Errorf("after encode/decode data changed!\nOriginal: %X\nOutput:   %X\n", data, x)
		}
	}
}
