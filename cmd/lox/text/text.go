package text

import (
	"tiny-script/interpreter"
	"tiny-script/lexer"
	"tiny-script/parser"
)

func Start(text string) {
	l := lexer.New(text)
	p := parser.New(l)
	statements, err := p.Parse()
	if err == nil && len(statements) != 0 {
		interpreter.Interpret(statements)
	}
}
