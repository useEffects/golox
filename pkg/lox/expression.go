package lox

type Expression struct {
	Left     *Expression
	Operator *Token
	Right    *Expression
}


