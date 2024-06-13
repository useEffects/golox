package main

import (
	"golox/pkg/lox"
	"log"
	"os"
)

func main() {
	lox := lox.Lox{}
	if len(os.Args) > 1 {
		log.Fatal("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		lox.RunFile(os.Args[1])
	} else {
		lox.RunPrompt()
	}
}
