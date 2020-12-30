package def

type Expr interface {
	AcceptStr(visitor StrVisitor) string
	Accept(interpreterVisitor InterpreterVisitor) interface{}
}

type EmptyExpr struct {
}

type Literal struct {
	Value interface{}
}

type Grouping struct {
	Expression Expr
}

type Binary struct {
	Left  Expr
	Token Token
	Right Expr
}

type Unary struct {
	Token Token
	Right Expr
}
