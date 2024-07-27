package lox

import "fmt"

type RuntimeError struct {
	Message string
	Token   Token
}

func (e *RuntimeError) Error() {
	panic(fmt.Sprintf("%s\n[line %d]", e.Message, e.Token.Line))
}
