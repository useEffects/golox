package main

import (
	"fmt"
	"golox/pkg/lox"
)

func main() {
	Test()
}

func CreateExpr() lox.Expr[string] {
	return &lox.Binary[string]{
		Left: &lox.Unary[string]{
			Operator: lox.Token{Type: lox.MINUS, Lexeme: "-"},
			Right:    &lox.Literal[string]{Value: 123},
		},
		Operator: lox.Token{Type: lox.STAR, Lexeme: "*"},
		Right:    &lox.Grouping[string]{Expression: &lox.Literal[string]{Value: 45.67}},
	}
}

func Test() {
	expr := CreateExpr()
	printer := lox.AstPrinter{}
	fmt.Println(printer.Print(expr))
}
