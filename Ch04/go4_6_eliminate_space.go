package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	arr := []byte("hello          world")
	eliminateSpace(&arr)
	fmt.Println(string(arr))
}

func eliminateSpace(arr* []byte) {
	if len(*arr) == 0 {
		return
	}
	
	writePos := 0
	prevIsSpace := false
	
	for readPos := 0; readPos < len(*arr); {
		r, size := utf8.DecodeRune((*arr)[readPos:])		// size is the size of each character, might be 1 to 4
		
		if unicode.IsSpace(r) {
			if !prevIsSpace {
				copy((*arr)[writePos:], (*arr)[readPos:readPos+size])
				writePos += size
				prevIsSpace = true
			}
		} else {
			copy((*arr)[writePos:], (*arr)[readPos:readPos+size])
			writePos += size
			prevIsSpace = false
		}
		
		readPos += size
	}
	*arr = (*arr)[:writePos]
}
