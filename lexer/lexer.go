package lexer

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"tiny-script/token"
)

const (
	EOF = rune(-1)
)

var (
	// identifer error
	errUnterminated = errors.New("unterminated string")
	errEspace       = errors.New("invalid escape char")
	errInvalidChar  = errors.New("invalid unicode char")

	// number error
	errLessPower = errors.New("power is required")
)

// Lexer represents a lexical scanner for Lox programing language.
type Lexer struct {
	str       []rune
	currIndex int
	ch        rune
	tokenBuf  *strings.Builder
}

func (l *Lexer) Next() rune {
	if l.currIndex >= len(l.str) {
		return EOF
	}
	ch := l.str[l.currIndex]
	l.currIndex++
	return ch
}

func (l *Lexer) consume() {
	if l.Eof() {
		return
	}
	ch := l.Next()
	if ch == EOF {
		l.ch = EOF
		return
	}
	l.ch = ch
}

func (l *Lexer) peek() rune {
	return l.str[l.currIndex]
}

func (l *Lexer) skip() {
	for unicode.IsSpace(l.ch) {
		l.consume()
	}
}

func (l *Lexer) Eof() bool {
	return l.ch == EOF
}

func (l *Lexer) match(ch rune) bool {
	l.consume()
	if l.Eof() || l.ch != ch {
		return false
	}
	l.consume()
	return true
}

func (l *Lexer) error(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "%d %s\n", l.Pos(), msg)
}

func (l *Lexer) readIdentifier() string {
	l.tokenBuf.Reset()
	for isAlphaNumeric(l.ch) {
		l.tokenBuf.WriteRune(l.ch)
		l.consume()
	}
	return l.tokenBuf.String()
}

func (l *Lexer) readString() (string, error) {
	l.tokenBuf.Reset()
	l.consume()
	if l.ch == '"' {
		l.consume()
		return "", nil
	}

	for l.ch != '"' {
		if l.Eof() {
			l.error(errUnterminated.Error())
			return "", errUnterminated
		} else if l.ch == '\\' {
			peekCh := l.peek()
			if peekCh == EOF {
				l.error(errEspace.Error())
				return "", errEspace
			}
			l.consume()
			switch peekCh {
			case '"':
				l.tokenBuf.WriteRune('"')
			case 'u':
				code := make([]rune, 4)
				for i := range code {
					l.consume()
					if !unicode.Is(unicode.Hex_Digit, l.ch) {
						l.error(errInvalidChar.Error())
						return "", errInvalidChar
					}
					code[i] = l.ch
				}
				l.tokenBuf.WriteRune(charCode2Rune(string(code)))
			}
		} else {
			l.tokenBuf.WriteRune(l.ch)
		}
		l.consume()
	}
	// end ".
	l.consume()
	return l.tokenBuf.String(), nil
}

func (l *Lexer) readNumber() (string, error) {

	l.tokenBuf.Reset()
	for unicode.IsNumber(l.ch) {
		l.tokenBuf.WriteRune(l.ch)
		l.consume()
	}

	if l.ch == '.' {
		if !unicode.IsNumber(l.peek()) {
			return l.tokenBuf.String(), nil
		}
		l.tokenBuf.WriteRune(l.ch)
		l.consume()
		for unicode.IsNumber(l.ch) {
			l.tokenBuf.WriteRune(l.ch)
			l.consume()
		}
	}

	if l.ch == 'E' || l.ch == 'e' {
		seenPower := false
		l.tokenBuf.WriteRune(l.ch)
		l.consume()
		if l.ch == '+' || l.ch == '-' {
			l.tokenBuf.WriteRune(l.ch)
			l.consume()
		}
		for unicode.IsNumber(l.ch) {
			seenPower = true
			l.tokenBuf.WriteRune(l.ch)
			l.consume()
		}
		if !seenPower {
			l.error(errLessPower.Error())
			return "", errLessPower
		}
	}

	return l.tokenBuf.String(), nil
}

// NextToken reads and returns token and literal.
// It returns token.Illegal for invalid string or number.
// It return token.EOF at the end of input string.
func (l *Lexer) NextToken() (tok token.Token, literal string) {
	l.skip()
	switch l.ch {
	case '&':
		tok = token.And
		literal = "&"
	case '|':
		tok = token.Or
		literal = "|"
	case '(':
		tok = token.LeftParen
		literal = "("
	case ')':
		tok = token.RightParen
		literal = ")"
	case '{':
		tok = token.LeftBrace
		literal = "{"
	case '}':
		tok = token.RightBrace
		literal = "}"
	case ',':
		tok = token.Comma
		literal = ","
	case '.':
		tok = token.Dot
		literal = "."
	case '-':
		tok = token.Minus
		literal = "-"
	case '+':
		tok = token.Plus
		literal = "+"
	case ';':
		tok = token.Semicolon
		literal = ";"
	case '/':
		tok = token.Slash
		literal = "/"
	case '*':
		tok = token.Star
		literal = "*"
	case '!':
		if l.match('=') {
			tok = token.BangEqual
			literal = "!="
		} else {
			tok = token.Bang
			literal = "!"
		}
		return
	case '=':
		if l.match('=') {
			tok = token.EqualEqual
			literal = "=="
		} else {
			tok = token.Equal
			literal = "="
		}
		return
	case '>':
		if l.match('=') {
			tok = token.GreaterEqual
			literal = ">="
		} else {
			tok = token.Greater
			literal = ">"
		}
		return
	case '<':
		if l.match('=') {
			tok = token.LessEqual
			literal = "<="
		} else {
			tok = token.Less
			literal = "<"
		}
		return
	case '"':
		liter, err := l.readString()
		if err != nil {
			return token.Illegal, liter
		}
		tok = token.String
		literal = liter
		return
	case EOF:
		tok = token.EOF
		return
	default:
		if unicode.IsLetter(l.ch) {
			literal = l.readIdentifier()
			tok = token.Lookup(literal)
			return
		} else if unicode.IsNumber(l.ch) {
			liter, err := l.readNumber()
			if err != nil {
				return token.Illegal, ""
			}
			tok = token.Number
			literal = liter
			return
		}

		tok = token.Illegal
		literal = ""
	}

	l.consume()
	return
}

// Pos returns current position of lexer.
func (l *Lexer) Pos() int {
	return l.currIndex
}

func charCode2Rune(code string) rune {
	v, err := strconv.ParseInt(code, 16, 32)
	if err != nil {
		return unicode.ReplacementChar
	}
	return rune(v)
}

func isAlphaNumeric(ch rune) bool { return unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == '_' }

// New return an instance of Lexer.
func New(input string) *Lexer {
	l := &Lexer{
		str:       []rune(input),
		currIndex: 0,
		tokenBuf:  &strings.Builder{},
	}
	l.consume()
	return l
}
