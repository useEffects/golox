package main

import (
	"golox/pkg/lox"
	"log"
	"os"
)

func main() {
	interpreter := lox.Interpreter{}
	if len(os.Args) > 1 {
		log.Fatal("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		interpreter.RunFile(os.Args[1])
	} else {
		interpreter.RunPrompt()
	}
}
