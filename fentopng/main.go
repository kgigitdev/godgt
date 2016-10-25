package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/kgigitdev/godgt"
)

// Quick and dirty fentopng command line utility.
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s fen [fen ...]\n", os.Args[0])
		fmt.Printf("cat fens.txt | %s -\n", os.Args[0])
		os.Exit(1)
	}
	if os.Args[1] == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		count := 1
		for scanner.Scan() {
			fen := scanner.Text()
			outfile := fmt.Sprintf("board-%03d.png", count)
			godgt.WritePng(fen, 64, outfile)
			count++
		}
	} else {
		for count, fen := range os.Args[1:] {
			outfile := fmt.Sprintf("board-%03d.png", count)
			godgt.WritePng(fen, 64, outfile)
			count++
		}
	}
}
