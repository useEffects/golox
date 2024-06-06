package main

import (
	"golox/pkg/lox"
	"log"
	"os"
)

func main() {
	interp := lox.Interpreter{}
	if len(os.Args) > 1 {
		log.Fatal("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		interp.RunFile(os.Args[1])
	} else {
		interp.RunPrompt()
	}
}
