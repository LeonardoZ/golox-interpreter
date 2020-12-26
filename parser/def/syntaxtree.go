package def

type StrVisitor interface {
	visitBinaryExprStr(binary *Binary) string
	visitUnaryExpStr(unary *Unary) string
	visitGroupingExprStr(grouping *Grouping) string
	visitLiteralExprStr(literal *Literal) string
}

type Expr interface {
	AcceptStr(visitor StrVisitor) string
}

type EmptyExpr struct {
}

func (empty *EmptyExpr) AcceptStr(v StrVisitor) string {
	return ""
}

type Literal struct {
	Value interface{}
}

func (literal *Literal) AcceptStr(v StrVisitor) string {
	return v.visitLiteralExprStr(literal)
}

type Grouping struct {
	Expression Expr
}

func (grouping *Grouping) AcceptStr(v StrVisitor) string {
	return v.visitGroupingExprStr(grouping)
}

type Binary struct {
	Left  Expr
	Token Token
	Right Expr
}

func (binary *Binary) AcceptStr(v StrVisitor) string {
	return v.visitBinaryExprStr(binary)
}

type Unary struct {
	Token Token
	Right Expr
}

func (unary *Unary) AcceptStr(v StrVisitor) string {
	return v.visitUnaryExpStr(unary)
}
