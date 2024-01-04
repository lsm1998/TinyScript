package main

import (
	"fmt"
	"os"
	"tiny-script/lox/repl"

	"tiny-script/interpreter"
	"tiny-script/lexer"
	"tiny-script/parser"
)

func main() {
	if len(os.Args) >= 2 {
		name := os.Args[1]
		b, err := os.ReadFile(name)
		if err != nil {
			panic(err)
		}
		l := lexer.New(string(b))
		p := parser.New(l)
		if statements, err := p.Parse(); err == nil && len(statements) != 0 {
			interpreter.Interpret(statements)
		}
		return
	}

	_, _ = fmt.Fprintln(os.Stdout, "TinyScript programing language.")
	_, _ = fmt.Fprintln(os.Stdout, "Feel free to type commands.")
	_, _ = fmt.Fprintln(os.Stdout, "Type \"exit\" to exit.")
	repl.Start(os.Stdin, os.Stdout)
}
