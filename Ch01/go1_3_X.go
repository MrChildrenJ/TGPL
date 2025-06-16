package main

import (
	"fmt"
	"os"
	"strings"
)


func main() {
	fmt.Println(strings.Join(os.Args[1:], " "))
	my_print(os.Args[1:])
}

func my_print([]string) {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

// Compare Join and my_print