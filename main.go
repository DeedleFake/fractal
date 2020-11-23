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
	px = -0.5557506
	py = -0.55560
	ph = 0.000000001
	//px = -2
	//py = -1.2
	//ph = 2.5

	// Quality
	imgWidth  = 1024
	imgHeight = 1024
	maxIter   = 1500
	samples   = 50

	profileCPU = true
)

const (
	ratio = float64(imgWidth) / float64(imgHeight)
)

func main() {
	fmt.Println("Allocating image...")
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

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
	wg.Add(imgHeight)
	for y := 0; y < imgHeight; y++ {
		go renderRow(&wg, img, y)
	}
	wg.Wait()

	showProgressDone()
}

func renderRow(wg *sync.WaitGroup, img *image.RGBA, y int) {
	defer wg.Done()
	defer showProgress(y)

	for x := 0; x < imgWidth; x++ {
		var r, g, b int
		for i := 0; i < samples; i++ {
			nx := ph*ratio*((float64(x)+RandFloat64())/float64(imgWidth)) + px
			ny := ph*((float64(y)+RandFloat64())/float64(imgHeight)) + py

			c := paint(mandelbrotIter(nx, ny, maxIter))
			r += mixColorPart(c.R)
			g += mixColorPart(c.G)
			b += mixColorPart(c.B)
		}

		cr := toRGBPart(float64(r) / float64(samples))
		cg := toRGBPart(float64(g) / float64(samples))
		cb := toRGBPart(float64(b) / float64(samples))

		setPix(img, x, y, color.RGBA{R: cr, G: cg, B: cb, A: 255})
	}
}

func setPix(p *image.RGBA, x, y int, c color.RGBA) {
	// Copied from (*image.RGBA).SetRGBA() to skip bounds check.
	i := y*p.Stride + x*4
	s := p.Pix[i : i+4 : i+4]
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
	s[3] = c.A
}

func paint(r float64, n int) color.RGBA {
	if r > 4 {
		return hslToRGB(float64(n)/800*r, 1, 0.5)
	}

	return color.RGBA{R: 255, G: 255, B: 255, A: 255}
}

func mandelbrotIter(px, py float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + px
		y = 2*xy + py
	}

	return xx + yy, maxIter
}

// by u/Boraini
//func mandelbrotIterComplex(px, py float64, maxIter int) (float64, int) {
//	var current complex128
//	pxpy := complex(px, py)
//
//	for i := 0; i < maxIter; i++ {
//		magnitude := cmplx.Abs(current)
//		if magnitude > 2 {
//			return magnitude * magnitude, i
//		}
//		current = current * current + pxpy
//	}
//
//	magnitude := cmplx.Abs(current)
//	return magnitude * magnitude, maxIter
//}
