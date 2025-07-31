package evaluator

import (
	"gogogo/lexer"
	"gogogo/object"
	"gogogo/parser"
)

// 测试用
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}
