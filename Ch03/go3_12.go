package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter 1st string: ")
	stringS, _ := reader.ReadString('\n')
	fmt.Print("Enter 2nd string: ")
	stringT, _ := reader.ReadString('\n')

	fmt.Println(compareStr(stringS, stringT))
}

func compareStr(s string, t string) bool {
	if len(s) != len(t) { 
		return false
	}

	ss := []byte(s)
	tt := []byte(t)

	sort.Slice(ss, func(i, j int) bool { return ss[i] < ss[j] })
	sort.Slice(tt, func(i, j int) bool { return tt[i] < tt[j] })

	for i := range ss {
		if ss[i] != tt[i] {
			return false
		}
	}
	return true
}