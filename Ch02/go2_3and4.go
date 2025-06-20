package main

import (
	"fmt"
	"time"
)

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountLoop(x uint64) int {
	var count int
	for i := 0; i < 8; i++ {
		count += int(pc[byte(x>>(i*8))])
	}
	return count
}

func PopCountOnce(x uint64) int {
	var count int
	for i := 0; i < 64; i++ {
		count += int(pc[byte(x>>i)])
	}
	return count
}

func BenchMark(name string, fn func(uint64) int, iterations int) {	// pass test func into BenchMark
	start := time.Now()

	testValues := []uint64{
		0x1234567890ABCDEF,
		0xFFFFFFFFFFFFFFFF,
		0x0000000000000000,
		0x8000000000000001,
	}

	for i := 0; i < iterations; i++ {
		for _, val := range testValues {
			fn(val)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("%s: %v\n", name, int(elapsed)/1000000)
}

func main() {
	// testValues := []uint64{
	// 	0x1234567890ABCDEF,
	// 	0xFFFFFFFFFFFFFFFF,
	// 	0x0000000000000000,
	// 	0x8000000000000001,
	// 	0x123867BC90A45DEF,
	// 	0xFFEFFFF7FFAFFFFF,
	// 	0x0010003005006000,
	// 	0x8001001000100001,
	// }

	// for _, val := range testValues {
	// 	origin := PopCount(val)
	// 	loop := PopCountLoop(val)
	// 	fmt.Printf("Origin: %v, Loop: %v\n", origin, loop)
	// }

	iterations := 10000000

	BenchMark("origin", PopCount, iterations)
	BenchMark("loop", PopCountLoop, iterations)
	BenchMark("once", PopCountOnce, iterations)
}