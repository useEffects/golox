package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var expressions = []string{
	"Binary => left:Expr, operator:Token, right:Expr",
	"Grouping => expression:Expr",
	"Literal => value:interface{}",
	"Unary => operator:Token, right:Expr",
}

func defineIntro(file *os.File) {
	intro := `
/* 
This is an auto-generated file by using the generate_ast.go utility program.
(っ◔◡◔)っ ♥ ast ♥
*/
`
	file.WriteString(intro)
}

func defineExports(file *os.File, types []string) {
	file.WriteString("package main\n\n")
	file.WriteString("import (\n\t\"fmt\"\n)\n\n")

	file.WriteString("type Visitor interface {\n")
	for _, t := range types {
		className := strings.TrimSpace(strings.Split(t, "=>")[0])
		file.WriteString(fmt.Sprintf("\tVisit%s(*%s) interface{}\n", className, className))
	}
	file.WriteString("}\n\n")

	file.WriteString("type Expr interface {\n")
	file.WriteString("\tAccept(Visitor) interface{}\n")
	file.WriteString("}\n\n")
}

func defineVisitor(file *os.File, className string, types []string) {
	file.WriteString("type Visitor interface {\n\n")
	for _, t := range types {
		typeName := strings.TrimSpace(strings.Split(t, "=>")[0])
		file.WriteString(fmt.Sprintf("\tVisit%s(*%s) interface{}\n", typeName, typeName))
	}
	file.WriteString("}\n\n")
}

func defineType(file *os.File, className string, baseClassName string, attributes []string) {
	file.WriteString(fmt.Sprintf("type %s struct {\n", className))
	for _, attr := range attributes {
		parts := strings.Split(strings.TrimSpace(attr), ":")
		if len(parts) != 2 {
			log.Fatalf("Invalid attribute format: %s", attr)
		}
		attrName := strings.TrimSpace(parts[0])
		attrType := strings.TrimSpace(parts[1])
		file.WriteString(fmt.Sprintf("\t%s %s\n", attrName, attrType))
	}
	file.WriteString("}\n\n")

	file.WriteString(fmt.Sprintf("func (t *%s) Accept(v Visitor) interface{} {\n", className))
	file.WriteString(fmt.Sprintf("\treturn v.Visit%s(t)\n", className))
	file.WriteString("}\n\n")
}

func generateAst(outputDir string, baseClassName string) {
	fileName := fmt.Sprintf("%s/%s.go", outputDir, strings.ToLower(baseClassName))
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Could not create file: %v", err)
	}
	defer file.Close()

	defineIntro(file)
	defineExports(file, expressions)
	defineVisitor(file, baseClassName, expressions)

	for _, item := range expressions {
		parts := strings.Split(item, "=>")
		if len(parts) != 2 {
			log.Fatalf("Invalid expression format: %s", item)
		}
		className := strings.TrimSpace(parts[0])
		attrs := strings.Split(strings.TrimSpace(parts[1]), ",")
		defineType(file, className, baseClassName, attrs)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Expected one argument: <output directory>")
		os.Exit(1)
	}
	outputDir := os.Args[1]
	fmt.Println("Generating AST to:", outputDir)
	generateAst(outputDir, "Expr")
}
