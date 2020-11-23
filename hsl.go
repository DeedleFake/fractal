package main

import "image/color"

// https://axonflux.com/handy-rgb-to-hsl-and-rgb-to-hsv-color-model-c

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	} else if t > 1 {
		t--
	}

	switch {
	case t < 1.0/6.0:
		return p + (q-p)*6*t
	case t < 1.0/2.0:
		return q
	case t < 2.0/3.0:
		return p + (q-p)*(2.0/3.0-t)*6
	default:
		return p
	}
}

func hslToRGB(h, s, l float64) color.RGBA {
	if s == 0 {
		return color.RGBA{R: uint8(l * 255), G: uint8(l * 255), B: uint8(l * 255), A: 255}
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	return color.RGBA{
		R: uint8(hueToRGB(p, q, h+1.0/3.0) * 255),
		G: uint8(hueToRGB(p, q, h) * 255),
		B: uint8(hueToRGB(p, q, h-1.0/3.0) * 255),
		A: 255,
	}
}
