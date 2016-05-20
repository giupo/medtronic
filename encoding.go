package medtronic

import (
	"errors"
)

func Encode4b6b(src []byte) []byte {
	// 2 input bytes produce 3 output bytes.
	// Odd final input byte, if any, produces 2 output bytes.
	n := len(src)
	dst := make([]byte, 3*(n/2)+2*(n%2))
	for i, j := 0, 0; i < n; i, j = i+2, j+3 {
		x := src[i]
		a := encode4b[hi(4, x)]
		b := encode4b[lo(4, x)]
		dst[j] = (a << 2) | hi(4, b)
		c := byte(0)
		d := byte(0)
		if i+1 < n {
			y := src[i+1]
			c = encode4b[hi(4, y)]
			d = encode4b[lo(4, y)]
			dst[j+2] = (lo(2, c) << 6) | d
		}
		dst[j+1] = (lo(4, b) << 4) | hi(6, c)
	}
	return dst
}

var DecodingFailure = errors.New("4b/6b decoding failure")

func Decode6b4b(src []byte) ([]byte, error) {
	n := len(src)
	// Check for valid packet length.
	if n%3 == 1 {
		return nil, DecodingFailure
	}
	// 3 input bytes produce 2 output bytes.
	// Final 2 input bytes, if any, produce 1 output byte.
	dst := make([]byte, 2*(n/3)+(n%3)/2)
	for i, j := 0, 0; i < n; i, j = i+3, j+2 {
		if i+1 >= n {
			return dst, DecodingFailure // shouldn't happen
		}
		x := src[i]
		y := src[i+1]
		a := decode6b[hi(6, x)]
		b := decode6b[(lo(2, x)<<4)|hi(4, y)]
		if a == 0xFF || b == 0xFF {
			return dst, DecodingFailure
		}
		dst[j] = (a << 4) | b
		if i+2 < n {
			z := src[i+2]
			c := decode6b[(lo(4, y)<<2)|hi(2, z)]
			d := decode6b[lo(6, z)]
			if c == 0xFF || d == 0xFF {
				return dst, DecodingFailure
			}
			dst[j+1] = (c << 4) | d
		}
	}
	return dst, nil
}

func hi(n, x byte) byte {
	return x >> (8 - n)
}

func lo(n, x byte) byte {
	return x & ((1 << n) - 1)
}

var (
	encode4b = []byte{
		0x15, 0x31, 0x32, 0x23,
		0x34, 0x25, 0x26, 0x16,
		0x1A, 0x19, 0x2A, 0x0B,
		0x2C, 0x0D, 0x0E, 0x1C,
	}

	// Inverse of encode4b table, with 0xFF indicating an undefined value.
	decode6b = []byte{
		/* 0x00 */ 0xFF /* 0x01 */, 0xFF /* 0x02 */, 0xFF /* 0x03 */, 0xFF,
		/* 0x04 */ 0xFF /* 0x05 */, 0xFF /* 0x06 */, 0xFF /* 0x07 */, 0xFF,
		/* 0x08 */ 0xFF /* 0x09 */, 0xFF /* 0x0A */, 0xFF /* 0x0B */, 0x0B,
		/* 0x0C */ 0xFF /* 0x0D */, 0x0D /* 0x0E */, 0x0E /* 0x0F */, 0xFF,
		/* 0x10 */ 0xFF /* 0x11 */, 0xFF /* 0x12 */, 0xFF /* 0x13 */, 0xFF,
		/* 0x14 */ 0xFF /* 0x15 */, 0x00 /* 0x16 */, 0x07 /* 0x17 */, 0xFF,
		/* 0x18 */ 0xFF /* 0x19 */, 0x09 /* 0x1A */, 0x08 /* 0x1B */, 0xFF,
		/* 0x1C */ 0x0F /* 0x1D */, 0xFF /* 0x1E */, 0xFF /* 0x1F */, 0xFF,
		/* 0x20 */ 0xFF /* 0x21 */, 0xFF /* 0x22 */, 0xFF /* 0x23 */, 0x03,
		/* 0x24 */ 0xFF /* 0x25 */, 0x05 /* 0x26 */, 0x06 /* 0x27 */, 0xFF,
		/* 0x28 */ 0xFF /* 0x29 */, 0xFF /* 0x2A */, 0x0A /* 0x2B */, 0xFF,
		/* 0x2C */ 0x0C /* 0x2D */, 0xFF /* 0x2E */, 0xFF /* 0x2F */, 0xFF,
		/* 0x30 */ 0xFF /* 0x31 */, 0x01 /* 0x32 */, 0x02 /* 0x33 */, 0xFF,
		/* 0x34 */ 0x04 /* 0x35 */, 0xFF /* 0x36 */, 0xFF /* 0x37 */, 0xFF,
		/* 0x38 */ 0xFF /* 0x39 */, 0xFF /* 0x3A */, 0xFF /* 0x3B */, 0xFF,
		/* 0x3C */ 0xFF /* 0x3D */, 0xFF /* 0x3E */, 0xFF /* 0x3F */, 0xFF,
	}
)