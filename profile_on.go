// +build profile

package main

import (
	"fmt"
	"os"
	"runtime/pprof"
)

func startProfiling() func() {
	f, err := os.Create("profile.prof")
	if err != nil {
		panic(fmt.Errorf("create CPU profile output file: %w", err))
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic(fmt.Errorf("start CPU profiling: %w", err))
	}

	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}
}
