package parser

type StrVisitor interface {
	visitBinaryExprStr(binary *Binary) string
	visitUnaryExpStr(unary *Unary) string
	visitGroupingExprStr(grouping *Grouping) string
	visitLiteralExprStr(literal *Literal) string
}

type Expr interface {
	acceptStr(visitor StrVisitor) string
}

type Literal struct {
	Value interface{}
}

func (literal *Literal) acceptStr(v StrVisitor) string {
	return v.visitLiteralExprStr(literal)
}

type Grouping struct {
	Expression Expr
}

func (grouping *Grouping) acceptStr(v StrVisitor) string {
	return v.visitGroupingExprStr(grouping)
}

type Binary struct {
	Left  Expr
	Token Token
	Right Expr
}

func (binary *Binary) acceptStr(v StrVisitor) string {
	return v.visitBinaryExprStr(binary)
}

type Unary struct {
	Token Token
	Right Expr
}

func (unary *Unary) acceptStr(v StrVisitor) string {
	return v.visitUnaryExpStr(unary)
}
