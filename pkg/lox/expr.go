package lox

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

// Binary represents a binary expression
type Binary[T any] struct {
	Left Expr[T]
	Operator Token
	Right Expr[T]
}

func (e Binary[T]) visit(visitor Visitor[T]) T {
	return visitor.VisitBinary(e)
}

// Grouping represents a grouping expression
type Grouping[T any] struct {
	Expression Expr[T]
}

func (e Grouping[T]) visit(visitor Visitor[T]) T {
	return visitor.VisitGrouping(e)
}

// Literal represents a literal expression
type Literal[T any] struct {
	Value interface{}
}

func (e Literal[T]) visit(visitor Visitor[T]) T {
	return visitor.VisitLiteral(e)
}

// Unary represents a unary expression
type Unary[T any] struct {
	Operator Token
	Right Expr[T]
}

func (e Unary[T]) visit(visitor Visitor[T]) T {
	return visitor.VisitUnary(e)
}

