package main

import (
	"math"
	"testing"
	"time"
)

func TestAnglePress(t *testing.T) {
	InitUART()
	for i := 0; i < 360; i++ {
		sendCommand(rstickAngle(int64(i), 0xFF))
		time.Sleep(1 * time.Millisecond)
	}
	sendNoInput()
}

func TestAngle(t *testing.T) {
	x := -10
	y := 0
	cVal := math.Atan(float64(y)/float64(x)) * 180 / math.Pi
	if x < 0 {
		cVal += 180
	} else if y < 0 {
		cVal += 360
	}
	t.Log(cVal)
}
