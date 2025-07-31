// ast/ast.go
package ast

import (
	"bytes"
	"fmt"
	"gogogo/token"
	"strings"
)

// 无论是语句还是标识符都会实现该方法
// 返回当前的token的标识,用于调试当前是否是对应的目标值
type Node interface {
	String() string
}

// 程序是由多个语句构成的
type Program struct {
	Statements []Statement
}

func (p Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// 语句接口,对下面的结构体进行分析
// 什么结构体可以是语句
// return、let、块语句{}、表达式语句
type Statement interface {
	Node
	// 返回当前语句的字符串表示
	statementNode()
	String() string
}

// x + 10;   // 可能有的语言并没有这样的
type ReturnStatement struct {
	Token       token.TokenType
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(string(rs.Token))
	out.WriteString(" ")
	out.WriteString(rs.ReturnValue.String())
	return out.String()
}

func (rs *ReturnStatement) statementNode() {}

// let x = 10;
type LetStatement struct {
	Token token.TokenType
	// 标识符
	Name *Identifier
	// 表达式
	Value Expression
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(string(ls.Token))
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	out.WriteString(ls.Value.String())
	return out.String()
}

func (ls *LetStatement) statementNode() {}

// 块语句
type BlockStatement struct {
	Token      token.TokenType
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() token.TokenType {
	return bs.Token
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// 表达式语句
type ExpressionStatement struct {
	Token      token.TokenType // 该表达式中的第一个词法单元
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	var out bytes.Buffer
	return es.Expression.String()
	return out.String()
}

func (es *ExpressionStatement) statementNode() {}

// 表达式接口,对下面的结构体进行分析
// 什么结构体可以是表达式
// identifier 、INT、还有就是前缀、中缀组合可以是表达式,还有call fn返回表达式结果,函数也是一个表达式,if 也是表达式
// 在本语言中if是一个表达式会返回值,if括号里的判断是由表达式构成,里面的内容是由语句构成(注意在本语言中表达式也是语句)
type Expression interface {
	Node
	expressionNode()
	TokenLiteral() token.TokenType
	// 返回当前表达式的字符串表示
}

// 标识符
type Identifier struct {
	Token token.TokenType // token.IDENT词法单元
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() token.TokenType {
	return i.Token
}

func (i *Identifier) String() string {
	return i.Value
}

type Boolean struct {
	Token token.TokenType
	Value bool
}

func (b *Boolean) TokenLiteral() token.TokenType {
	return b.Token
}

func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) expressionNode() {}

type StringLiteral struct {
	Token token.TokenType
	Value string
}

func (b *StringLiteral) TokenLiteral() token.TokenType {
	return b.Token
}

func (b *StringLiteral) String() string {
	return fmt.Sprintf("%s", b.Value)
}

func (b *StringLiteral) expressionNode() {}

type IntegerLiteral struct {
	Token token.TokenType
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() token.TokenType {
	return il.Token
}

func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

func (il *IntegerLiteral) expressionNode() {}

// if
type IfExpression struct {
	Token token.TokenType // 'if'词法单元
	// 条件
	Condition Expression
	// then
	Consequence *BlockStatement
	// else
	Alternative *BlockStatement
}

func (ie *IfExpression) TokenLiteral() token.TokenType {
	return ie.Token
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString("(")
	out.WriteString(ie.Condition.String())
	out.WriteString(") ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

func (ie *IfExpression) expressionNode() {}

type FunctionLiteral struct {
	Token token.TokenType // 'fn'词法单元
	// 参数
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() token.TokenType {
	return fl.Token
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(string(fl.TokenLiteral()))
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.TokenType // '('词法单元
	Function  Expression      // 标识符或函数字面量
	Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() token.TokenType {
	return ce.Token
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(")")
	return out.String()
}

func (ce *CallExpression) expressionNode() {}

// 前缀表达式是由符号和另一个表达式构成的
type PrefixExpression struct {
	Token    token.TokenType // 前缀词法单元，如!
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() token.TokenType {
	return pe.Token
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

func (pe *PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Token    token.TokenType // 运算符词法单元，如+
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() token.TokenType {
	return ie.Token
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

func (ie *InfixExpression) expressionNode() {}

// 数组是由多个表达式构成的
type ArrayLiteral struct {
	Token    token.TokenType // token.IDENT词法单元
	Elements []Expression
}

func (ar *ArrayLiteral) TokenLiteral() token.TokenType {
	return ar.Token
}

func (ar *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ar.Elements {
		elements = append(elements, e.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (ar *ArrayLiteral) expressionNode() {}

// 用于查找数组下标的情况
type IndexExpression struct {
	Token token.TokenType // '['词法单元
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) TokenLiteral() token.TokenType {
	return ie.Token
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

func (ie *IndexExpression) expressionNode() {}
