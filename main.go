package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

// Configuration
const (
	// Position and height
	Position = -0.5557506 - 0.55560i
	Height   = 0.000000001
	//Position = -2 - 1.2i
	//Height = 2.5

	// Quality
	ImageWidth    = 1024.0
	ImageHeight   = 1024.0
	MaxIterations = 1500
	Samples       = 50
	Threshold     = 4
)

const (
	ratio = ImageWidth / ImageHeight
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
		panic(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done!")
}

// render renders a fractal to img.
func render(img *image.RGBA) {
	if profileCPU {
		f, err := os.Create("profile.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var wg sync.WaitGroup
	wg.Add(ImageHeight)
	for y := 0; y < ImageHeight; y++ {
		go renderRow(&wg, img, y)
	}
	wg.Wait()

	showProgressDone()
}

// renderRow renders a single row of pixels to img, telling wg when
// it's done.
func renderRow(wg *sync.WaitGroup, img *image.RGBA, y int) {
	defer wg.Done()
	defer showProgress(y)

	fy := float64(y)
	for x := 0; x < ImageWidth; x++ {
		fx := float64(x)
		var r, g, b int
		for i := 0; i < Samples; i++ {
			c := Height*complex(
				ratio*((fx+randFloat64())/ImageWidth),
				(fy+randFloat64())/ImageHeight,
			) + Position

			col := mandelbrotColor(mandelbrotIter(c))
			r += colorStep(col.R)
			g += colorStep(col.G)
			b += colorStep(col.B)
		}

		cr := convertColor(float64(r) / float64(Samples))
		cg := convertColor(float64(g) / float64(Samples))
		cb := convertColor(float64(b) / float64(Samples))

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
		return hslToRGB(float64(iter)/800*check, 1, 0.5)
	}

	return color.RGBA{R: 255, G: 255, B: 255, A: 255}
}

// mandelbrotIter checks if |f(z)| becomes greater than a threshold
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
