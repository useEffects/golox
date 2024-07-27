package scanner

import (
	"fmt"
	"strconv"

	"golox/pkg/fault"
)

type scanner struct {
	Source  string
	Tokens  []Token
	start   int
	current int
	line    int
	err     error
}

func NewScanner(source string) *scanner {
	tokens := make([]Token, 0, 10)
	return &scanner{source, tokens, 0, 0, 1, nil}
}

func (s *scanner) ScanTokens() error {
	for s.current < len(s.Source) {
		s.start = s.current
		switch s.Source[s.current] {
		case '(':
			s.addToken(LEFT_PAREN, nil)
		case ')':
			s.addToken(RIGHT_PAREN, nil)
		case '{':
			s.addToken(LEFT_BRACE, nil)
		case '}':
			s.addToken(RIGHT_BRACE, nil)
		case ',':
			s.addToken(COMMA, nil)
		case '.':
			s.addToken(DOT, nil)
		case '-':
			s.addToken(MINUS, nil)
		case '+':
			s.addToken(PLUS, nil)
		case ';':
			s.addToken(SEMICOLON, nil)
		case '*':
			s.addToken(STAR, nil)
		case '!':
			if s.next('=') {
				s.addToken(BANG_EQUAL, nil)
			} else {
				s.addToken(BANG, nil)
			}
		case '=':
			if s.next('=') {
				s.addToken(EQUAL_EQUAL, nil)
			} else {
				s.addToken(EQUAL, nil)
			}
		case '<':
			if s.next('=') {
				s.addToken(LESS_EQUAL, nil)
			} else {
				s.addToken(LESS, nil)
			}
		case '>':
			if s.next('=') {
				s.addToken(GREATER_EQUAL, nil)
			} else {
				s.addToken(GREATER, nil)
			}
		case '/':
			if s.next('/') {
				s.singleComment()
			} else {
				s.addToken(SLASH, nil)
			}
		case ' ':
		case '\t':
		case '\r':
		case '\n':
			s.line++
		case '"':
			err := s.string()
			if s.err == nil {
				s.err = err
			}
		default:
			if isDigit(s.Source[s.current]) {
				s.number()
			} else if isAlpha(s.Source[s.current]) {
				s.identifier()
			} else {
				message := fmt.Sprintf("unknown character '%c'", s.Source[s.current])
				s.err = fault.NewFault(s.line, message)
			}
		}
		s.current++
	}
	s.Tokens = append(s.Tokens, Token{EOF, "EOF", nil, s.line})
	return s.err
}

func (s *scanner) singleComment() {
	for s.current < len(s.Source) && s.Source[s.current] != '\n' {
		s.current++
	}
	s.current--
}

func (s *scanner) string() error {
	s.current++
	for s.current < len(s.Source) && s.Source[s.current] != '"' {
		if s.Source[s.current] == '\n' {
			s.line++
		}
		s.current++
	}

	if s.current == len(s.Source) {
		return fault.NewFault(s.line, "unterminated string")
	} else {
		s.addToken(STRING, s.Source[s.start+1:s.current])
	}

	return nil
}

func (s *scanner) number() {
	for s.current < len(s.Source) && isDigit(s.Source[s.current]) {
		s.current++
	}

	if s.current+1 < len(s.Source) && s.Source[s.current] == '.' && isDigit(s.Source[s.current+1]) {
		s.current++
		for s.current < len(s.Source) && isDigit(s.Source[s.current]) {
			s.current++
		}
	}

	s.current--
	value, err := strconv.ParseFloat(s.Source[s.start:s.current+1], 64)
	if err == nil {
		s.addToken(NUMBER, value)
	}
}

func (s *scanner) identifier() {
	for s.current < len(s.Source) && (isAlpha(s.Source[s.current]) || isDigit(s.Source[s.current])) {
		s.current++
	}

	s.current--
	lexeme := s.Source[s.start : s.current+1]
	if keyword, ok := keywords[lexeme]; ok {
		s.addToken(keyword, nil)
	} else {
		s.addToken(IDENTIFIER, nil)
	}
}

func (s *scanner) addToken(tokenType int, literal interface{}) {
	lexeme := s.Source[s.start : s.current+1]
	token := Token{tokenType, lexeme, literal, s.line}
	s.Tokens = append(s.Tokens, token)
}

func (s *scanner) next(c byte) bool {
	if s.current < len(s.Source)-1 && s.Source[s.current+1] != c {
		return false
	}

	s.current++
	return true
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}