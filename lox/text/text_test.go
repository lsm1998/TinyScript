package text

import (
	"strings"
	"testing"
	"text/scanner"
	"tiny-script/interpreter"
	"tiny-script/lexer"
	"tiny-script/parser"
)

var script = `
function gen() {
    let a = 0;
    function inner() {
        a = a + 1;
        return a;
    }
    return inner;
}

if (1 > 0 & 1 > 1)  {
	print "true";
} else {
	print "false";
}

let fn = gen();

print fn();
print fn();
`

func TestStart(t *testing.T) {
	Start(script)
}

func TestLexer(t *testing.T) {
	s := &scanner.Scanner{}
	s.Init(strings.NewReader("01234"))
	t.Log(string(s.Next()))

	l := lexer.New(script)

	p := parser.New(l)

	stmts, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

	for _, v := range stmts {
		t.Log(v.String())
	}
	interpreter.Interpret(stmts)
}
