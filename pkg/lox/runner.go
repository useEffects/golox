package lox

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Interpreter struct {
	HadError bool
}

func (i *Interpreter) RunFile(path string) {
	bytes, _ := os.ReadFile(path)

	i.Run(string(bytes))
}

func (i *Interpreter) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == os.ErrClosed {
				break
			}
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		line = line[:len(line)-1]
		i.Run(line)
	}
}

func (i *Interpreter) Run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token.ToString())
	}
}

func (i *Interpreter) Error(line int, message string) {
	i.Report(line, "", message)
}

func (i *Interpreter) Report(line int, where, message string) {
	log.Fatalln("[line", line, "] Error", where, ":", message)
	i.HadError = true
}
