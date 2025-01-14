package scanner

type Scanner struct {
	text string

	startI   int
	currentI int

	line   int
	column int
}

func NewScanner(text string) *Scanner {
	return &Scanner{
		text:     text,
		startI:   0,
		currentI: 0,
		line:     1,
	}
}

func (s *Scanner) Next() Token {
	s.skipWhitespace()
	s.startI = s.currentI

	if s.isAtEnd() {
		return s.make(Eof)
	}

	c := s.advance()

	if isAlpha(c) {
		return s.identifier()
	}
	if isDigit(c) || (c == '-' && isDigit(s.peek())) {
		return s.number(c)
	}

	switch c {
	case '(':
		return s.make(LeftParen)
	case ')':
		return s.make(RightParen)
	case '{':
		return s.make(LeftBrace)
	case '}':
		return s.make(RightBrace)
	case '[':
		return s.make(LeftBracket)
	case ']':
		return s.make(RightBracket)

	case '.':
		return s.make(Dot)
	case ',':
		return s.make(Comma)
	case ':':
		return s.make(Colon)
	case ';':
		return s.make(Semicolon)

	case '+':
		if s.match('+') {
			return s.make(PlusPlus)
		}

		return s.matchToken('=', PlusEqual, Plus)
	case '-':
		if s.match('-') {
			return s.make(MinusMinus)
		}

		return s.matchToken('=', MinusEqual, Minus)
	case '*':
		return s.matchToken('=', StarEqual, Star)
	case '/':
		return s.matchToken('=', SlashEqual, Slash)
	case '%':
		return s.matchToken('=', PercentageEqual, Percentage)

	case '=':
		if s.match('=') {
			return s.make(EqualEqual)
		}

		return s.matchToken('>', FuncPtr, Equal)
	case '!':
		return s.matchToken('=', BangEqual, Bang)
	case '<':
		if s.match('<') {
			return s.matchToken('=', LessLessEqual, LessLess)
		}

		return s.matchToken('=', LessEqual, Less)
	case '>':
		if s.match('>') {
			return s.matchToken('=', GreaterGreaterEqual, GreaterGreater)
		}

		return s.matchToken('=', GreaterEqual, Greater)

	case '|':
		if s.match('=') {
			return s.make(PipeEqual)
		}

		return s.matchToken('|', Or, Pipe)
	case '^':
		if s.match('=') {
			return s.make(XorEqual)
		}

		return s.make(Xor)
	case '&':
		if s.match('=') {
			return s.make(AmpersandEqual)
		}

		return s.matchToken('&', And, Ampersand)

	case '#':
		return s.make(Hashtag)

	case '\'':
		return s.character()
	case '"':
		return s.string()
	}

	return s.error("Unexpected character.")
}

func (s *Scanner) identifier() Token {
	for isAlpha(s.peek()) || isDigit(s.peek()) {
		s.advance()
	}

	return s.make(s.identifierKind())
}

func (s *Scanner) identifierKind() TokenKind {
	switch s.text[s.startI] {
	case 'a':
		return s.checkKeyword(1, "s", As)
	case 'b':
		return s.checkKeyword(1, "reak", Break)
	case 'c':
		return s.checkKeyword(1, "ontinue", Continue)
	case 'e':
		if s.currentI-s.startI > 1 {
			switch s.text[s.startI+1] {
			case 'l':
				return s.checkKeyword(2, "se", Else)
			case 'n':
				return s.checkKeyword(2, "um", Enum)
			}
		}
	case 'f':
		if s.currentI-s.startI > 1 {
			switch s.text[s.startI+1] {
			case 'a':
				return s.checkKeyword(2, "lse", False)
			case 'o':
				return s.checkKeyword(2, "r", For)
			case 'u':
				return s.checkKeyword(2, "nc", Func)
			}
		}
	case 'i':
		if s.currentI-s.startI > 1 {
			switch s.text[s.startI+1] {
			case 'f':
				return If
			case 'm':
				return s.checkKeyword(2, "pl", Impl)
			}
		}
	case 'n':
		return s.checkKeyword(1, "il", Nil)
	case 'r':
		return s.checkKeyword(1, "eturn", Return)
	case 's':
		if s.currentI-s.startI > 1 {
			switch s.text[s.startI+1] {
			case 't':
				if s.currentI-s.startI > 2 {
					switch s.text[s.startI+2] {
					case 'a':
						return s.checkKeyword(3, "tic", Static)
					case 'r':
						return s.checkKeyword(3, "uct", Struct)
					}
				}
			}
		}
	case 't':
		return s.checkKeyword(1, "rue", True)
	case 'v':
		return s.checkKeyword(1, "ar", Var)
	case 'w':
		return s.checkKeyword(1, "hile", While)
	}

	return Identifier
}

func (s *Scanner) checkKeyword(start int, rest string, kind TokenKind) TokenKind {
	if s.currentI-s.startI == start+len(rest) && s.text[s.startI+start:s.startI+start+len(rest)] == rest {
		return kind
	}

	return Identifier
}

func (s *Scanner) number(c uint8) Token {
	next := s.peek()

	// Hex
	if c == '0' && (next == 'x' || next == 'X') {
		s.advance()
		return s.hex()
	}

	// Binary
	if c == '0' && (next == 'b' || next == 'B') {
		s.advance()
		return s.binary()
	}

	// Integers or floats
	return s.integerOrFloat()
}

func (s *Scanner) integerOrFloat() Token {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	if s.peek() == 'f' || s.peek() == 'F' {
		s.advance()
	}

	return s.make(Number)
}

func (s *Scanner) hex() Token {
	for isHex(s.peek()) {
		s.advance()
	}

	return s.make(Hex)
}

func (s *Scanner) binary() Token {
	for isBinary(s.peek()) {
		s.advance()
	}

	return s.make(Binary)
}

func (s *Scanner) string() Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		return s.error("Unterminated string")
	}

	s.advance()
	return s.make(String)
}

func (s *Scanner) character() Token {
	if s.isAtEnd() || s.peek() == '\'' {
		return s.error("Empty character.")
	}

	if s.advance() == '\\' && !s.isAtEnd() {
		c := s.advance()

		if c != '\'' && c != '0' && c != 'n' && c != 'r' && c != 't' {
			return s.error("Unexpected character.")
		}
	}

	if s.peek() != '\'' {
		return s.error("Unterminated character.")
	}

	s.advance()
	return s.make(Character)
}

func (s *Scanner) matchToken(expected uint8, kindTrue TokenKind, kindFalse TokenKind) Token {
	if s.match(expected) {
		return s.make(kindTrue)
	}

	return s.make(kindFalse)
}

func (s *Scanner) match(expected uint8) bool {
	if s.isAtEnd() {
		return false
	}

	if s.peek() != expected {
		return false
	}

	s.advance()
	return true
}

func (s *Scanner) skipWhitespace() {
	for {
		if s.isAtEnd() {
			return
		}

		c := s.peek()

		switch c {
		case ' ', '\r', '\t':
			s.advance()

		case '\n':
			s.advance()
			s.line++
			s.column = 0

		case '/':
			if s.peekNext() == '/' {
				for !s.isAtEnd() && s.peek() != '\n' {
					s.advance()
				}
			} else if s.peekNext() == '*' {
				s.advance()
				s.advance()

				for !s.isAtEnd() && (s.peek() != '*' || s.peekNext() != '/') {
					if s.peek() == '\n' {
						s.line++
						s.column = 0
					}

					s.advance()
				}

				if !s.isAtEnd() {
					s.advance()

					if !s.isAtEnd() {
						s.advance()
					}
				}
			} else {
				return
			}

		default:
			return
		}
	}
}

func (s *Scanner) peek() uint8 {
	if s.isAtEnd() {
		return '\000'
	}

	return s.text[s.currentI]
}

func (s *Scanner) peekNext() uint8 {
	if s.isAtEnd() {
		return '\000'
	}

	return s.text[s.currentI+1]
}

func (s *Scanner) advance() uint8 {
	s.currentI++
	s.column++

	return s.text[s.currentI-1]
}

func (s *Scanner) isAtEnd() bool {
	return s.currentI >= len(s.text)
}

func (s *Scanner) make(kind TokenKind) Token {
	lexeme := s.text[s.startI:s.currentI]

	return Token{
		Kind:   kind,
		Lexeme: lexeme,
		line:   s.line,
		column: s.column - len(lexeme),
	}
}

func (s *Scanner) error(msg string) Token {
	return Token{
		Kind:   Error,
		Lexeme: msg,
		line:   s.line,
		column: s.column,
	}
}

func isAlpha(c uint8) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}

func isHex(c uint8) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isBinary(c uint8) bool {
	return c == '0' || c == '1'
}
