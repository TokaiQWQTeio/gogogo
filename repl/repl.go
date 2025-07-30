package repl

import (
	"bufio"
	"fmt"
	"gogogo/evaluator"
	"gogogo/lexer"
	"gogogo/object"
	"gogogo/parser"
	"gogogo/token"
	"io"
	"strings"
)

const (
	RESET  = "\033[0m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PURPLE = "\033[35m"
	CYAN   = "\033[36m"
	WHITE  = "\033[37m"
)

const PROMPT = ">> "

func colorizeInput(input string) string {
	l := lexer.New(input)
	var colored strings.Builder

	for {
		tok := l.NextToken()
		if tok.Type == token.EOF {
			break
		}

		switch tok.Type {
		case token.LET, token.FUNCTION, token.IF, token.ELSE, token.TRUE, token.FALSE, token.RETURN:
			colored.WriteString(GREEN + tok.Literal + RESET)
		case token.INT:
			colored.WriteString(BLUE + tok.Literal + RESET)
		case token.STRING:
			colored.WriteString(PURPLE + tok.Literal + RESET)
		case token.ASSIGN, token.PLUS, token.MINUS, token.BANG, token.ASTERISK, token.SLASH:
			colored.WriteString(RED + tok.Literal + RESET)
		case token.EQ, token.NOT_EQ, token.LT, token.GT:
			colored.WriteString(RED + tok.Literal + RESET)
		case token.IDENT:
			colored.WriteString(WHITE + tok.Literal + RESET)
		default:
			colored.WriteString(tok.Literal)
		}
	}

	return colored.String()
}

//func Start(in io.Reader, out io.Writer) {
//	scanner := bufio.NewScanner(in)
//	env := object.NewEnvironment()
//	for {
//		fmt.Fprintf(out, PROMPT)
//		scanned := scanner.Scan()
//		if !scanned {
//			return
//		}
//		line := scanner.Text()
//		l := lexer.New(line)
//		p := parser.New(l)
//
//		program := p.ParseProgram()
//		if len(p.Errors()) != 0 {
//			printParserErrors(out, p.Errors())
//			continue
//		}
//
//		evaluated := evaluator.Eval(program, env)
//		if evaluated != nil {
//			io.WriteString(out, evaluated.Inspect())
//			io.WriteString(out, "\n")
//		}
//	}
//}

// 修改 Start 函数
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment() // 如果不需要立即执行可以注释掉

	fmt.Fprintf(out, "Monkey REPL with syntax highlighting\n")
	fmt.Fprintf(out, "Press Ctrl+C to exit\n")

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// 显示带颜色的输入
		coloredLine := colorizeInput(line)
		fmt.Fprintf(out, "Colored: %s\n", coloredLine)

		// 继续正常的解析和执行流程
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// 如果需要执行，取消下面的注释
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`
