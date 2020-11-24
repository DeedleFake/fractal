// +build !profile

package main

const profileCPU = false

func startProfiling() func() {
	return func() {}
}
