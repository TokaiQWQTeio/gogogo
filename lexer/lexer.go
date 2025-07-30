package lexer

import "gogogo/token"

type Lexer struct {
	input    string // 输入的字符串进行此法分析的
	offset   int    // 当前对于字符串的偏移
	rdOffset int    // 预读取的位置
	ch       byte   // offset 偏移的位置的字符
}

// 创建一个lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// 每一次读取一个tokenType
// 注意对于上面的一些case并没有return操作,下面的一些有return操作
// 没有return会再次执行readChar()然后在返回
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		//判断是不是==
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '!':
		//判断是不是!=
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// 读取标识符、数字、或者什么都不是的ILLEGAL
		if isLetter(l.ch) {
			tok.Literal = l.readTarget(isLetter)
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			//是不是数字
			tok.Type = token.INT
			tok.Literal = l.readTarget(isDigit)
			return tok
		} else {
			//什么都不是
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// 如果是一些类似于空格的特殊字符进行读取下一个操作
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// rdOffset 代表的是读取的下一个的位置，如果当前的位置已经到达了输入的长度，代表已经到达了EOF
// 到达了EOF 将l.ch = 0,否则将l.ch 设置为l.rdOffset的位置的值
// 然后让offset和rdOffset向后移动
func (l *Lexer) readChar() {
	if l.rdOffset >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.rdOffset]
	}
	l.offset = l.rdOffset
	l.rdOffset += 1
}

// 查看下一个byte的值是什么,并不进行偏移设置
func (l *Lexer) peekChar() byte {
	if l.rdOffset >= len(l.input) {
		return 0
	} else {
		return l.input[l.rdOffset]
	}
}

// 如果遇到属于标识符的范围的字符进行继续读取操作,直到l.ch不是范围内的字符
// 根据切片的左闭右开原则,进行字符串的截取操作,传入的参数可以是标识符函数和数字函数
func (l *Lexer) readTarget(f func(ch byte) bool) string {
	// 记录当前位置
	position := l.offset
	for f(l.ch) {
		l.readChar()
	}
	// 返回截取的字符串
	return l.input[position:l.offset]
}

// 跳过最初的" 以及最后的 "
func (l *Lexer) readString() string {
	position := l.offset + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.offset]
}

// 判断是不是字母
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// 判断是不是引号
func isQuotation(ch byte) bool {
	return ch == '"'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
