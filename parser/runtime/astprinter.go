package runtime

import (
	"fmt"
	"loxlang/parser/def"
)

// AstPrinter - implements Visitor Pattern
type AstPrinter struct {
}

// Print Prints the result of the AST
func (astPrinter *AstPrinter) Print(expr def.Expr) {
	fmt.Println(">: " + expr.AcceptStr(astPrinter))
}

// VisitBinaryExprStr Handles Binary
func (astPrinter *AstPrinter) VisitBinaryExprStr(binary *def.Binary) string {
	return astPrinter.parenthesize(binary.Token.Lexeme, binary.Left, binary.Right)
}

// VisitUnaryExpStr Handles Unary
func (astPrinter *AstPrinter) VisitUnaryExpStr(unary *def.Unary) string {
	return astPrinter.parenthesize(unary.Token.Lexeme, unary.Right)
}

// VisitGroupingExprStr Handles Grouping
func (astPrinter *AstPrinter) VisitGroupingExprStr(grouping *def.Grouping) string {
	return astPrinter.parenthesize("group", grouping.Expression)
}

// VisitLiteralExprStr Handles Literal
func (astPrinter *AstPrinter) VisitLiteralExprStr(literal *def.Literal) string {
	if literal.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", literal.Value)
}

// VisitVariableExprStr Handles Variable
func (astPrinter *AstPrinter) VisitVariableExprStr(variable *def.Variable) string {
	return fmt.Sprintf("%s", variable.Name.Literal)
}

func (astPrinter *AstPrinter) parenthesize(name string, exprs ...def.Expr) string {
	var result string
	result += "(" + name
	for _, e := range exprs {
		result += " "
		result += e.AcceptStr(astPrinter)
	}
	result += ")"
	return result

}
