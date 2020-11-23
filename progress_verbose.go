// +build !quiet

package main

import "fmt"

func showProgress(y int) {
	fmt.Printf("\r%d/%d (%d%%)", y, imgHeight, int(100*(float64(y)/float64(imgHeight))))
}

func showProgressDone() {
	fmt.Printf("\r%d/%[1]d (100%%)\n", imgHeight)
}
