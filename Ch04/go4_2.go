package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)

func main() {
	// flag.String(name, defaultValue, usage)
	hashType := flag.String("hash", "256", "Hash algorithm: 256, 384, or 512")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please enter input: ")
	if !scanner.Scan() {
		fmt.Fprintf(os.Stderr, "Error reading input\n")
		os.Exit(1)
	}

	input := scanner.Text()					// input is string
	// input, err := io.ReadAll(os.Stdin)	// input is []byte
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	// 	os.Exit(1)
	// }

	switch *hashType {
	case "256":
		hash := sha256.Sum256([]byte(input))	// need to convert input to []byte
		fmt.Printf("%x\n", hash)
	case "384":
		hash := sha512.Sum384([]byte(input))
		fmt.Printf("%x\n", hash)
	case "512":
		hash := sha512.Sum512([]byte(input))
		fmt.Printf("%x\n", hash)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported hash type: %s\n", *hashType)
		os.Exit(1)
	}	
}