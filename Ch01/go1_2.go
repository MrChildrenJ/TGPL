package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main2() {
	for i, arg := range os.Args {
		if i == 0 {
			arg = filepath.Base(arg)
		}
		fmt.Println(i, arg)
	}
}