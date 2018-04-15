package main

import (
	"fmt"
	"os"

	"github.com/waits/tang/file"
	"github.com/waits/tang/repl"
)

func main() {
	switch len(os.Args) {
	case 1:
		fmt.Printf("Tang 0.1.0-alpha\n")
		repl.Start(os.Stdin, os.Stdout)
	case 2:
		file.Exec(os.Stdout, os.Args[1], os.Args[2:])
	default:
		panic("too many arguments")
	}
}
