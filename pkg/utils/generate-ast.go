package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		os.Stderr.WriteString("Usage: generate_ast <output directory>\n")
		os.Exit(64)
	}
	outDir := os.Args[1]
	defineAst(outDir, "Expr", []string{
		"Binary   : Left Expr[T], Operator Token, Right Expr[T]",
		"Grouping : Expression Expr[T]",
		"Literal  : Value interface{}",
		"Unary    : Operator Token, Right Expr[T]",
	})
}

func defineAst(outputDir string, baseName string, types []string) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file %s\n", path)
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString("package lox\n\n")
	file.WriteString("// Expr represents an expression in the AST\n")
	file.WriteString("type Expr[T any] interface {\n")
	file.WriteString("\tvisit(visitor Visitor[T]) T\n")
	file.WriteString("}\n\n")
	defineVisitor(file, types)
	file.WriteString("\n")
	for _, t := range types {
		structName := strings.TrimSpace(t[:strings.Index(t, ":")])
		fields := strings.TrimSpace(t[strings.Index(t, ":")+1:])
		defineType(file, structName, fields)
	}
}

func defineType(file *os.File, structName string, fieldList string) {
	file.WriteString("// " + structName + " represents a " + strings.ToLower(structName) + " expression\n")
	file.WriteString("type " + structName + "[T any] struct {\n")
	fields := strings.Split(fieldList, ", ")
	for _, field := range fields {
		file.WriteString("\t" + field + "\n")
	}
	file.WriteString("}\n\n")

	// Implement visit method for the struct
	file.WriteString("func (e " + structName + "[T]) visit(visitor Visitor[T]) T {\n")
	file.WriteString("\treturn visitor.Visit" + structName + "(e)\n")
	file.WriteString("}\n\n")
}

func defineVisitor(file *os.File, types []string) {
	file.WriteString("// Visitor defines the visitor interface for expression traversal\n")
	file.WriteString("type Visitor[T any] interface {\n")
	for _, t := range types {
		typeName := strings.TrimSpace(t[:strings.Index(t, ":")])
		file.WriteString("\tVisit" + typeName + "(" + typeName + "[T]) T\n")
	}
	file.WriteString("}\n")
}
