package parser

import (
	"encoding/json"
	"fmt"
	"gogogo/lexer"
	"testing"
)

func TestString(t *testing.T) {
	input := `let add = fn(x, y) {
  x + y;
};
let result = add(five, ten);
if (5 < 10) {
	return true;
} else {
	return false;
}
let ten = 10;
let five = 5+10*30;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	s, _ := json.MarshalIndent(program, "", "  ")
	fmt.Println(string(s))
	checkParserErrors(t, p)
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
