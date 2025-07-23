package lexer

type TokenType int

const (
	TOKEN_ILLEGAL TokenType = iota
	TOKEN_EOF

	TOKEN_IDENT   // variable names, block types
	TOKEN_STRING  // "value"
	TOKEN_BOOLEAN // true, false
	TOKEN_INT     // 123
	TOKEN_LONG    // 1234567890
	TOKEN_FLOAT   // 123.456
	TOKEN_DOUBLE  // 123.456789 (higher precision)

	TOKEN_ASSIGN // =
	TOKEN_COLON  // :
	TOKEN_COMMA  // ,
	TOKEN_DOT    // .
	TOKEN_LPAREN // (
	TOKEN_RPAREN // )
	TOKEN_LBRACE // {
	TOKEN_RBRACE // }

	TOKEN_KEYWORD // keywords: frame, var, slot, trigger, etc.
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type Lexer struct {
	input        string
	position     int  // current char index
	readPosition int  // next char index
	ch           byte // current char
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l._readChar()
	return l
}

func (l *Lexer) NextToken() Token {
	l._skipWhitespace()

	// Skip comments
	if l.ch == '/' && l._peekChar() == '/' {
		l._skipComment()
		return l.NextToken()
	}

	switch l.ch {
	case '=':
		return l._newToken(TOKEN_ASSIGN, string(l.ch))
	case ':':
		return l._newToken(TOKEN_COLON, string(l.ch))
	case ',':
		return l._newToken(TOKEN_COMMA, string(l.ch))
	case '.':
		return l._newToken(TOKEN_DOT, string(l.ch))
	case '(':
		return l._newToken(TOKEN_LPAREN, string(l.ch))
	case ')':
		return l._newToken(TOKEN_RPAREN, string(l.ch))
	case '{':
		return l._newToken(TOKEN_LBRACE, string(l.ch))
	case '}':
		return l._newToken(TOKEN_RBRACE, string(l.ch))
	case '"':
		return l._readString()
	case 0:
		return l._newToken(TOKEN_EOF, "")
	default:
		if _isLetter(l.ch) {
			literal := l._readIdentifier()
			tokenType := _lookupKeyword(literal)
			return Token{
				Type:    tokenType,
				Literal: literal,
				Line:    l.line,
				Column:  l.column,
			}
		} else if _isDigit(l.ch) || l.ch == '.' {
			return l._readNumber()
		}
		tok := Token{
			Type:    TOKEN_ILLEGAL,
			Literal: string(l.ch),
			Line:    l.line,
			Column:  l.column,
		}
		l._readChar()
		return tok
	}
}

func (l *Lexer) _readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) _peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) _skipComment() {
	// Skip the '//' characters
	l._readChar()
	l._readChar()

	// Continue reading until the end of line or EOF
	for l.ch != '\n' && l.ch != 0 {
		l._readChar()
	}
}

func (l *Lexer) _newToken(tokenType TokenType, ch string) Token {
	tok := Token{
		Type:    tokenType,
		Literal: ch,
		Line:    l.line,
		Column:  l.column,
	}
	l._readChar()
	return tok
}

func (l *Lexer) _skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l._readChar()
	}
}

func (l *Lexer) _readIdentifier() string {
	start := l.position
	for _isLetter(l.ch) || _isDigit(l.ch) || l.ch == '_' {
		l._readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) _readString() Token {
	startLine, startCol := l.line, l.column
	l._readChar() // skip initial quote
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l._readChar()
	}
	literal := l.input[start:l.position]
	l._readChar() // skip closing quote
	return Token{
		Type:    TOKEN_STRING,
		Literal: literal,
		Line:    startLine,
		Column:  startCol,
	}
}

func (l *Lexer) _readNumber() Token {
	startLine, startCol := l.line, l.column
	start := l.position
	seenDot := false

	for _isDigit(l.ch) || (!seenDot && l.ch == '.') {
		if l.ch == '.' {
			seenDot = true
		}
		l._readChar()
	}

	literal := l.input[start:l.position]
	numberType := TOKEN_INT

	if seenDot {
		dotIdx := -1
		for i, c := range literal {
			if c == '.' {
				dotIdx = i
				break
			}
		}
		fracLen := 0
		if dotIdx != -1 {
			fracLen = len(literal) - dotIdx - 1
		}
		if fracLen > 6 {
			numberType = TOKEN_DOUBLE
		} else {
			numberType = TOKEN_FLOAT
		}
	} else {
		if len(literal) > 9 {
			numberType = TOKEN_LONG
		} else {
			numberType = TOKEN_INT
		}
	}

	return Token{
		Type:    numberType,
		Literal: literal,
		Line:    startLine,
		Column:  startCol,
	}
}

func _lookupKeyword(lit string) TokenType {
	switch lit {
	case "frame", "var", "slot", "trigger", "prop", "action", "block", "then", "props", "data":
		return TOKEN_KEYWORD
	case "true", "false":
		return TOKEN_BOOLEAN
	default:
		return TOKEN_IDENT
	}
}

func _isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func _isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
