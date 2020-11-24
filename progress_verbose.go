// +build !quiet

package main

import "fmt"

var (
	progressChan     = make(chan struct{})
	progressDoneChan = make(chan struct{})
)

func init() {
	go func() {
		defer close(progressDoneChan)
		defer fmt.Println()

		var done int
		for range progressChan {
			done++
			fmt.Printf("\r%v/%v (%v%%)", done, ImageHeight, int(float64(done)*100/ImageHeight))
		}
	}()
}

func updateProgress() {
	progressChan <- struct{}{}
}

func progressDone() {
	close(progressChan)
	<-progressDoneChan
}
