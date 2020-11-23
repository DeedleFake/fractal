// +build !quiet

package main

import "fmt"

func showProgress(y int) {
	fmt.Printf("\r%v/%v (%v%%)", y, ImageHeight, int(float64(y)*100/ImageHeight))
}

func showProgressDone() {
	fmt.Printf("\r%v/%[1]v (100%%)\n", ImageHeight)
}
