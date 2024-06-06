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
		"Binary   : Left Expr, Operator Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value interface{}",
		"Unary    : Operator Token, Right Expr",
	})
}

func defineAst(outputDir string, baseName string, types []string) {
	path := outputDir + "/" + baseName + ".go"
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open file %s\n", path)
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString("package " + baseName + "\n\n")
	file.WriteString("type " + baseName + " struct {\n")
	defineVisitor(file, baseName, types)
	file.WriteString("}\n\n")
	for _, t := range types {
		structName := t[:strings.Index(t, ":")]
		fields := t[strings.Index(t, ":")+1:]
		defineType(file, structName, fields)
	}
	file.WriteString("\n")
	file.WriteString("func (v " + baseName + ".Visitor) Accept(expr Expr) interface{} {}\n")
}

func defineType(file *os.File, structName string, fieldList string) {
	file.WriteString("type " + structName + " struct {\n")
	fields := strings.Split(fieldList, ", ")
	for _, field := range fields {
		file.WriteString("\t" + field + "\n")
	}
	file.WriteString("}\n")
}

func defineVisitor(file *os.File, baseName string, types []string) {
	file.WriteString("\tVisitor struct {\n")
	for _, t := range types {
		typeName := t[:strings.Index(t, ":")]
		file.WriteString("\tVisit" + typeName + " func(" + typeName + ") interface{}\n")
	}
	file.WriteString("\t}\n")
}
