package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Configuration constants.
const (
	Position = -0.5557506 - 0.55560i
	Height   = 0.000000001
	//Position = -2 - 1.2i
	//Height   = 2.5

	ImageWidth    = 1280
	ImageHeight   = 1024
	MaxIterations = 1500
	Samples       = 50
	Threshold     = 4

	IterHueAdjust = 800
)

// Configuration variables.
var (
	ThresholdColor = color.RGBAModel.Convert(color.White).(color.RGBA)
)

func main() {
	fmt.Println("Allocating image...")
	img := image.NewRGBA(image.Rect(0, 0, ImageWidth, ImageHeight))

	fmt.Println("Rendering...")
	start := time.Now()
	render(img)
	end := time.Now()

	fmt.Println("Done rendering in", end.Sub(start))

	fmt.Println("Encoding image...")
	f, err := os.Create("result.png")
	if err != nil {
		panic(fmt.Errorf("create image output file: %w", err))
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		panic(fmt.Errorf("encode image: %w", err))
	}
	fmt.Println("Done!")
}

// render renders a fractal to img.
func render(img *image.RGBA) {
	defer startProfiling()()
	defer progressDone()

	var wg sync.WaitGroup
	wg.Add(ImageHeight)
	for y := 0; y < ImageHeight; y++ {
		go renderRow(&wg, img, y)
	}
	wg.Wait()
}

// renderRow renders a single row of pixels to img, telling wg when
// it's done.
func renderRow(wg *sync.WaitGroup, img *image.RGBA, y int) {
	defer wg.Done()
	defer updateProgress(1)

	s := rand.New(rand.NewSource(int64(time.Now().UnixNano())))

	for x := 0; x < ImageWidth; x++ {
		xy := complex(float64(x), float64(y))

		var r, g, b int
		for i := float64(0); i < Samples; i++ {
			shifted := xy + complex(s.Float64(), s.Float64())
			c := Height*complex(real(shifted)/ImageWidth, imag(shifted)/ImageHeight) + Position

			col := mandelbrotColor(mandelbrotIter(c))
			r += colorStep(col.R)
			g += colorStep(col.G)
			b += colorStep(col.B)
		}

		cr := convertColor(float64(r) / Samples)
		cg := convertColor(float64(g) / Samples)
		cb := convertColor(float64(b) / Samples)

		setPix(img, x, y, color.RGBA{R: cr, G: cg, B: cb, A: 255})
	}
}

// setPix sets a pixel in an image.RGBA to a given color. It's
// basically directly copied from (*image.RGBA).SetRGBA() to skip the
// bounds check.
func setPix(p *image.RGBA, x, y int, c color.RGBA) {
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4]
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
	s[3] = c.A
}

// mandelbrotColor returns the color that a pixel should be colored
// based on the results of the Mandelbrot iteration for that pixel.
func mandelbrotColor(check float64, iter int) color.RGBA {
	if check > Threshold {
		return hslToRGB(float64(iter)/IterHueAdjust*check, 1, 0.5)
	}

	return ThresholdColor
}

// mandelbrotIter checks if |f(z)| becomes greater than Threshold
// when repeatedly applied to its own output, starting from z = 0,
// where f(z) = z*z + c.
//
// It returns |f(z)|^2 for the final result of f(z) and the number of
// times that it iterated to get to that result.
func mandelbrotIter(c complex128) (check float64, iter int) {
	prev := c
	for ; (iter < MaxIterations) && (check <= Threshold); iter++ {
		check = real(prev)*real(prev) + imag(prev)*imag(prev)
		prev = prev*prev + c
	}
	return check, iter
}
