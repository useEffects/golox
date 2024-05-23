package main

import (
	"golox/interpreter"
	"log"
	"os"
)

func main() {
	interp := interpreter.Interpreter{}
	if len(os.Args) > 1 {
		log.Fatal("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		interp.RunFile(os.Args[1])
	} else {
		interp.RunPrompt()
	}
}
