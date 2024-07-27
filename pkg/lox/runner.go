package lox

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Lox struct {
	HadError        bool
	HadRuntimeError bool
	interpreter     Interpreter[string]
}

func (i *Lox) RunFile(path string) {
	bytes, _ := os.ReadFile(path)

	i.Run(string(bytes))

	if i.HadError {
		os.Exit(65)
	}
	if i.HadRuntimeError {
		os.Exit(70)
	}
}

func (i *Lox) RunPrompt() {
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

func (i *Lox) Run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := Parser[string]{tokens: tokens}
	expression := parser.Parse()

	if i.HadError {
		return
	}

	// fmt.Println(AstPrinter{}.Print(expression))
	i.interpreter.Interpret(expression)
}

func (i *Lox) Error(param any, message string) {
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

func (i *Lox) RuntimeError(error *RuntimeError) {
	fmt.Fprintln(os.Stderr, error.Message, "\n[line", error.Token.Line, "]")
	i.HadRuntimeError = true
}

func (i *Lox) Report(line int, where, message string) {
	log.Fatalln("[line", line, "] Error", where, ":", message)
	i.HadError = true
}

var lox = Lox{}
