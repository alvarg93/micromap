package color

import (
	"encoding/hex"
	"math"
	"math/rand"
)

type RGB struct {
	R uint8
	G uint8
	B uint8
}

type HSV struct {
	H float64
	S float64
	V float64
}

type Color struct {
	RGB
	HSV
}

func NewRGB(r uint8, g uint8, b uint8) RGB {
	return RGB{
		R: r,
		G: g,
		B: b,
	}
}

func NewHSV(h float64, s float64, v float64) HSV {
	return HSV{
		H: h,
		S: s,
		V: v,
	}
}

func (hsv *HSV) ToRGB() RGB {
	hsv.S = math.Max(math.Min(hsv.S, 1.0), 0)
	hsv.V = math.Max(math.Min(hsv.V, 1.0), 0)
	C := hsv.V * hsv.S
	H := hsv.H / 60
	Hint := int(H)
	X := C * (1.0 - math.Abs(float64(Hint%2-Hint+-1)+H))
	m := hsv.V - C
	var r, g, b float64
	if H >= 0 && H < 1 {
		r = C
		g = X
		b = 0
	} else if H >= 1 && H < 2 {
		r = X
		g = C
		b = 0
	} else if H >= 2 && H < 3 {
		r = 0
		g = C
		b = X
	} else if H >= 3 && H < 4 {
		r = 0
		g = X
		b = C
	} else if H >= 4 && H < 5 {
		r = X
		g = 0
		b = C
	} else if H >= 5 && H < 6 {
		r = C
		g = 0
		b = X
	}

	newR := uint8(math.Max(math.Min((r+m)*255, 255), 0))
	newG := uint8(math.Max(math.Min((g+m)*255, 255), 0))
	newB := uint8(math.Max(math.Min((b+m)*255, 255), 0))
	return NewRGB(newR, newG, newB)
}

func RandomHSV(s, v float64) HSV {
	v = math.Max(math.Min(v, 1), 0)
	s = math.Max(math.Min(s, 1), 0)
	hsv := HSV{H: float64(rand.Uint32() % 360), S: s, V: v}
	return hsv
}

func RandomRGB() RGB {
	return NewRGB(uint8(rand.Uint32()%127+100), uint8(rand.Uint32()%127+100), uint8(rand.Uint32()%127+100))
}

func (rgb RGB) ToColorHash() string {
	return "#" + hex.EncodeToString([]byte{byte(rgb.R), byte(rgb.G), byte(rgb.B)})
}
