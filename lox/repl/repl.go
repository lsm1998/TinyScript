package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"tiny-script/interpreter"
	"tiny-script/lexer"
	"tiny-script/parser"
)

const prompt = ">> "

// Start creates a REPL for Lox.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	interpreter.SetEvalEnv("repl")
	for {
		fmt.Fprintf(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if line == "exit" {
			fmt.Fprintln(out, "bye.")
			return
		}
		l := lexer.New(line)
		p := parser.New(l)
		statements, err := p.Parse()
		if err == nil && len(statements) != 0 {
			interpreter.Interpret(statements)
		}
	}
}
