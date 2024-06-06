package lox

import (
	"fmt"
	"strings"
)

// Expr represents an expression in the AST
type Expr[T any] interface {
	visit(visitor Visitor[T]) T
}

// Visitor defines the visitor interface for expression traversal
type Visitor[T any] interface {
	VisitBinary(Binary[T]) T
	VisitGrouping(Grouping[T]) T
	VisitLiteral(Literal[T]) T
	VisitUnary(Unary[T]) T
}

// Binary represents a binary expression (left operand, operator, right operand)
type Binary[T any] struct {
	Left     Expr[T] // Left operand
	Operator Token   // Operator token
	Right    Expr[T] // Right operand
}

func (b Binary[T]) visit(visitor Visitor[T]) T {
	return visitor.VisitBinary(b)
}

// Grouping represents a grouped expression
type Grouping[T any] struct {
	Expression Expr[T] // Inner expression
}

func (g Grouping[T]) visit(visitor Visitor[T]) T {
	return visitor.VisitGrouping(g)
}

// Literal represents a literal value
type Literal[T any] struct {
	Value interface{} // Literal value (e.g., int, float)
}

func (l Literal[T]) visit(visitor Visitor[T]) T {
	return visitor.VisitLiteral(l)
}

// Unary represents a unary expression (operator, operand)
type Unary[T any] struct {
	Operator Token   // Operator token
	Right    Expr[T] // Operand
}

func (u Unary[T]) visit(visitor Visitor[T]) T {
	return visitor.VisitUnary(u)
}

// AstPrinter can be used to print the expression in a human-readable format
type AstPrinter struct{}

func (p AstPrinter) print(expr Expr[string]) string {
	return expr.visit(p)
}

func (p AstPrinter) VisitBinary(expr Binary[string]) string {
	return parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p AstPrinter) VisitGrouping(expr Grouping[string]) string {
	return parenthesize("group", expr.Expression)
}

func (p AstPrinter) VisitLiteral(expr Literal[string]) string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p AstPrinter) VisitUnary(expr Unary[string]) string {
	return parenthesize(expr.Operator.Lexeme, expr.Right)
}

func parenthesize(name string, exprs ...Expr[string]) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.visit(&AstPrinter{}))
	}
	builder.WriteString(")")
	return builder.String()
}

// CreateExpr constructs an expression based on the provided structure
func CreateExpr() Expr[string] {
	return &Binary[string]{
		Left: &Unary[string]{
			Operator: Token{Type: MINUS, Lexeme: "-"},
			Right:    &Literal[string]{Value: 123},
		},
		Operator: Token{Type: STAR, Lexeme: "*"},
		Right:    &Grouping[string]{Expression: &Literal[string]{Value: 45.67}},
	}
}

func Test() {
	expr := CreateExpr()
	printer := AstPrinter{}
	fmt.Println(printer.print(expr))
}
