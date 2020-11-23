// +build !nonlinear

package main

func mixColorPart(c uint8) int {
	return int(RGBToLinear(c))
}

func toRGBPart(c float64) uint8 {
	return LinearToRGB(uint64(c))
}
