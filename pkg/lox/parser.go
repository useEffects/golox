package lox

type Parser[T any] struct {
	tokens  []Token
	current int
}

type ParseError struct {
	message string
}

func (e *ParseError) Error() string {
	return e.message
}

func (p *Parser[T]) expression() Expr[T] {
	return p.equality()
}

func (p *Parser[T]) equality() Expr[T] {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary[T]{expr, operator, right}
	}

	return expr
}

func (p *Parser[T]) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser[T]) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser[T]) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser[T]) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser[T]) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser[T]) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser[T]) comparison() Expr[T] {
	expr := p.term()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Binary[T]{expr, operator, right}
	}
	return expr
}

func (p *Parser[T]) term() Expr[T] {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary[T]{expr, operator, right}
	}
	return expr
}

func (p *Parser[T]) factor() Expr[T] {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary[T]{expr, operator, right}
	}
	return expr
}

func (p *Parser[T]) unary() Expr[T] {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary[T]{operator, right}
	}
	expr, _ := p.primary()
	return expr
}

func (p *Parser[T]) primary() (Expr[T], ParseError) {
	if p.match(FALSE) {
		return Literal[T]{false}, ParseError{}
	}
	if p.match(TRUE) {
		return Literal[T]{true}, ParseError{}
	}
	if p.match(NIL) {
		return Literal[T]{nil}, ParseError{}
	}
	if p.match(NUMBER, STRING) {
		return Literal[T]{p.previous().Literal}, ParseError{}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping[T]{expr}, ParseError{}
	}

	return nil, p.Error(p.peek(), "Expect expression.")
}

func (p *Parser[T]) consume(t TokenType, message string) Token {
	if p.check(t) {
		return p.advance()
	}
	panic([]interface{}{p.peek(), message})
}

func (p *Parser[T]) Error(token Token, message string) ParseError {
	interpreter.Error(token.Line, message)
	return ParseError{}
}

func (p *Parser[T]) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser[T]) Parse() Expr[T] {
	return p.expression()
}
