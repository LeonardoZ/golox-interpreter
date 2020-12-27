package def

import (
	"fmt"
)

// AstPrinter - implements Visitor Pattern
type AstPrinter struct {
}

func (empty *EmptyExpr) AcceptStr(v StrVisitor) string {
	return ""
}

func (literal *Literal) AcceptStr(v StrVisitor) string {
	return v.visitLiteralExprStr(literal)
}

func (grouping *Grouping) AcceptStr(v StrVisitor) string {
	return v.visitGroupingExprStr(grouping)
}

func (binary *Binary) AcceptStr(v StrVisitor) string {
	return v.visitBinaryExprStr(binary)
}

func (unary *Unary) AcceptStr(v StrVisitor) string {
	return v.visitUnaryExpStr(unary)
}

type StrVisitor interface {
	visitBinaryExprStr(binary *Binary) string
	visitUnaryExpStr(unary *Unary) string
	visitGroupingExprStr(grouping *Grouping) string
	visitLiteralExprStr(literal *Literal) string
}

func (astPrinter *AstPrinter) Print(expr Expr) {
	fmt.Println(expr.AcceptStr(astPrinter))
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
		result += e.AcceptStr(astPrinter)
	}
	result += ")"
	return result

}
