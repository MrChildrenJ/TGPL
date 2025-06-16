package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main1() {
	s, sep := "", ""
	for i, arg := range os.Args {
		if i == 0 {
			arg = filepath.Base(arg)
		}
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}