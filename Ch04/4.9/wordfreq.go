package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type WordCount struct {
	Word string
	Count int
}

func main() {
	file, err := os.Open("TheGreatGatsby.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	count := make(map[string]int)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		count[word]++
	}

	// for word, count := range count {
	// 	println(word, count)
	// }

	// If want to print by order of frequency
	var wcs []WordCount
	for word, count := range count {
		wcs = append(wcs, WordCount{word, count})
	}

	sort.Slice(wcs, func(i, j int) bool {
		return wcs[i].Count > wcs[j].Count
	})

	for i, wc := range wcs {
		if i > 10 {
			break
		}
		fmt.Printf("%s: %d\n", wc.Word, wc.Count)
	}
}