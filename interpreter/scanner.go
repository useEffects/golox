package interpreter

var interp = Interpreter{}

type Scanner struct {
	source   string
	tokens   []Token
	start    int
	current  int
	line     int
	keywords map[string]TokenType
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		s.line++
	case '"':
		s.string()
	case 'o':
		if s.match('r') {
			s.addToken(OR)
		}
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			interp.Error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", nil, s.current})
	return s.tokens
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
		keywords: map[string]TokenType{
			"and":    AND,
			"class":  CLASS,
			"else":   ELSE,
			"false":  FALSE,
			"for":    FOR,
			"fun":    FUN,
			"if":     IF,
			"nil":    NIL,
			"or":     OR,
			"print":  PRINT,
			"return": RETURN,
			"super":  SUPER,
			"this":   THIS,
			"true":   TRUE,
			"var":    VAR,
			"while":  WHILE,
		},
	}
}

func (s *Scanner) advance() rune {
	s.current++
	return rune(s.source[s.current-1])
}

func (s *Scanner) addToken(tokenType TokenType, args ...interface{}) {
	text := s.source[s.start:s.current]
	var literal interface{}
	if len(args) > 0 {
		literal = args[0]
	}
	s.tokens = append(s.tokens, Token{tokenType, text, literal, s.line})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) string() {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		interp.Error(s.line, "Unterminated string.")
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	s.addToken(NUMBER, s.source[s.start:s.current])
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := s.keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}

	s.addToken(tokenType)
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}
