package main

import (
	"time"
)

// xorshift random

var randState = uint64(time.Now().UnixNano())

func randUint64() uint64 {
	randState = ((randState ^ (randState << 13)) ^ (randState >> 7)) ^ (randState << 17)
	return randState
}

func randFloat64() float64 {
	return float64(randUint64()/2) / (1 << 63)
}
