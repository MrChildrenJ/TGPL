package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)


func main3() {
	start := time.Now()
	fmt.Fprint(io.Discard, strings.Join(os.Args[1:], " "))
	secs1 := time.Since(start).Microseconds()

	start2 := time.Now()
	my_print(os.Args[1:])
	secs2 := time.Since(start2).Microseconds()

	fmt.Printf("Join took %d microseconds, my_print took %d microseconds\n", secs1, secs2)
}

func my_print([]string) {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Fprint(io.Discard, s)
}

// Compare Join and concatenation performance.
// For 500 words, Join took 0.000065 seconds, my_print took 0.000259 seconds