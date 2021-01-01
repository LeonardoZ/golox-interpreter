package def

import (
	"fmt"
)

// StrVisitor Interface
type StrVisitor interface {
	visitBinaryExprStr(binary *Binary) string
	visitUnaryExpStr(unary *Unary) string
	visitGroupingExprStr(grouping *Grouping) string
	visitLiteralExprStr(literal *Literal) string
	visitVariableExprStr(variable *Variable) string
}

// AstPrinter - implements Visitor Pattern
type AstPrinter struct {
}

func (empty *EmptyExpr) acceptStr(v StrVisitor) string {
	return ""
}

func (literal *Literal) acceptStr(v StrVisitor) string {
	return v.visitLiteralExprStr(literal)
}

func (grouping *Grouping) acceptStr(v StrVisitor) string {
	return v.visitGroupingExprStr(grouping)
}

func (binary *Binary) acceptStr(v StrVisitor) string {
	return v.visitBinaryExprStr(binary)
}

func (unary *Unary) acceptStr(v StrVisitor) string {
	return v.visitUnaryExpStr(unary)
}

func (variable *Variable) acceptStr(v StrVisitor) string {
	return v.visitVariableExprStr(variable)
}

// Print Prints the result of the AST
func (astPrinter *AstPrinter) Print(expr Expr) {
	fmt.Println(">: " + expr.acceptStr(astPrinter))
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

func (astPrinter *AstPrinter) visitVariableExprStr(variable *Variable) string {
	return fmt.Sprintf("%s", variable.Name.Literal)
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
