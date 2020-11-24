// +build !quiet

package main

import "fmt"

var (
	progressChan     = make(chan int)
	progressDoneChan = make(chan struct{})
)

func init() {
	go func() {
		defer close(progressDoneChan)
		defer fmt.Println()

		var done int
		for amount := range progressChan {
			done += amount
			fmt.Printf("\r%v/%v (%v%%)", done, ImageHeight, int(float64(done)*100/ImageHeight))
		}
	}()
}

func updateProgress(amount int) {
	progressChan <- amount
}

func progressDone() {
	close(progressChan)
	<-progressDoneChan
}
