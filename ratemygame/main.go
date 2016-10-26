package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

func main() {
	var opts Opts
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.ParseArgs(os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	g := NewGameRater(opts)
	g.Run()
}
