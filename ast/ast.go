// ast/ast.go
package ast

import (
	"gogogo/token"
)

// 无论是语句还是标识符都会实现该方法
// 返回当前的token的标识,用于调试当前是否是对应的目标值
type Node interface{}

// 程序是由多个语句构成的
type Program struct {
	Statements []Statement
}

// 语句接口,对下面的结构体进行分析
// 什么结构体可以是语句
// return、let、块语句{}、表达式语句
type Statement interface {
	Node
	statementNode()
}

// x + 10;   // 可能有的语言并没有这样的
type ReturnStatement struct {
	Token       token.TokenType
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// let x = 10;
type LetStatement struct {
	Token token.TokenType
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

type BlockStatement struct {
	Token      token.TokenType
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

type ExpressionStatement struct {
	Token      token.TokenType // 该表达式中的第一个词法单元
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// 表达式接口,对下面的结构体进行分析
// 什么结构体可以是表达式
// identifier 、INT、还有就是前缀、中缀组合可以是表达式,还有call fn返回表达式结果,函数也是一个表达式,if 也是表达式
// 在本语言中if是一个表达式会返回值,if括号里的判断是由表达式构成,里面的内容是由语句构成(注意在本语言中表达式也是语句)
type Expression interface {
	Node
	expressionNode()
	// 返回当前表达式的字符串表示
}

type Identifier struct {
	Token token.TokenType // token.IDENT词法单元
	Value string
}

func (i *Identifier) expressionNode() {}

type Boolean struct {
	Token token.TokenType
	Value bool
}

func (b *Boolean) expressionNode() {}

type StringLiteral struct {
	Token token.TokenType
	Value string
}

func (b *StringLiteral) expressionNode() {}

type IntegerLiteral struct {
	Token token.TokenType
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

type IfExpression struct {
	Token       token.TokenType // 'if'词法单元
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

type FunctionLiteral struct {
	Token      token.TokenType // 'fn'词法单元
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

type CallExpression struct {
	Token     token.TokenType // '('词法单元
	Function  Expression      // 标识符或函数字面量
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// 前缀表达式是由符号和另一个表达式构成的
type PrefixExpression struct {
	Token    token.TokenType // 前缀词法单元，如!
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Token    token.TokenType // 运算符词法单元，如+
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// 数组是由多个表达式构成的
type ArrayLiteral struct {
	Token    token.TokenType // token.IDENT词法单元
	Elements []Expression
}

func (ar *ArrayLiteral) expressionNode() {}

// 用于查找数组下标的情况
type IndexExpression struct {
	Token token.TokenType // '['词法单元
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
