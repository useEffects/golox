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

	parser := Parser[string]{tokens: tokens}
	expression := parser.Parse()

	if i.HadError {
		return
	}

	fmt.Println(AstPrinter{}.Print(expression))
}

func (i *Interpreter) Error(param any, message string) {
	switch p := param.(type) {
	case int:
		i.Report(p, "", message)
	case Token:
		if p.Type == EOF {
			i.Report(p.Line, " at end", message)
		} else {
			i.Report(p.Line, " at '"+p.Lexeme+"'", message)
		}
	default:
		panic("Invalid parameter type")
	}
}

func (i *Interpreter) Report(line int, where, message string) {
	log.Fatalln("[line", line, "] Error", where, ":", message)
	i.HadError = true
}

var interpreter = Interpreter{}
