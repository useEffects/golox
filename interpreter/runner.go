package interpreter

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Interpreter struct {
	HadError bool
}

func (i *Interpreter) RunFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return i.Run(string(bytes))
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
		if err := i.Run(line); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (i *Interpreter) Run(source string) error {
	return nil
}

func (i *Interpreter) Error(line int, message string) {
	i.Report(line, "", message)
}

func (i *Interpreter) Report(line int, where, message string) {
	log.Fatalln("[line", line, "] Error", where, ":", message)
	i.HadError = true
}
