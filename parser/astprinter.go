package parser

import "fmt"

// AstPrinter - implements Visitor Pattern
type AstPrinter struct {
}

func (astPrinter *AstPrinter) Print(expr Expr) {
	fmt.Println(expr.acceptStr(astPrinter))
}

func (astPrinter *AstPrinter) visitBinaryExprStr(binary *Binary) string {
	return astPrinter.parenthesize(binary.Token.Lexeme, binary.Left, binary.Right)
}

func (astPrinter *AstPrinter) visitUnaryExpStr(unary *Unary) string {
	return astPrinter.parenthesize(unary.Token.Lexeme, unary.Right)
}

func (astPrinter *AstPrinter) visitGroupingExprStr(grouping *Grouping) string {
	return astPrinter.parenthesize("group", grouping.Expression)
}

func (astPrinter *AstPrinter) visitLiteralExprStr(literal *Literal) string {
	if literal.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", literal.Value)
}

func (astPrinter *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var result string
	result += "(" + name
	for _, e := range exprs {
		result += " "
		result += e.acceptStr(astPrinter)
	}
	result += ")"
	return result

}
