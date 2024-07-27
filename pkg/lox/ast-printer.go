package lox

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p AstPrinter) Print(expr Expr[string]) string {
	return expr.accept(p)
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
		builder.WriteString(expr.accept(&AstPrinter{}))
	}
	builder.WriteString(")")
	return builder.String()
}
