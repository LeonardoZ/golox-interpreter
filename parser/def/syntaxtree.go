package def

// Stmt Base type for two kind of statements: Expression and Print. More to come.
type Stmt interface {
	accept(visitor StatementVisitor) *RuntimeError
}

// ExprStmt Expression Statements
type ExprStmt struct {
	Expr Expr
}

// Var Variable Declaration
type Var struct {
	Name        Token
	Initializer Expr
}

// Print is a simple Print statement for the language
type Print struct {
	Expr Expr
}

// Expr Mostly generic Tree Node
type Expr interface {
	acceptStr(visitor StrVisitor) string
	accept(interpreterVisitor ExpressionVisitor) (interface{}, *RuntimeError)
}

// EmptyExpr Just an "empty value" implementation for Expr
type EmptyExpr struct {
}

// Literal represents literal values like "abc", 13, 15.6
type Literal struct {
	Value interface{}
}

// Grouping represents parenthesised expressions
type Grouping struct {
	Expression Expr
}

// Binary represents expressions with two expr and one operator, like 1 + 2, a > b
type Binary struct {
	Left  Expr
	Token Token
	Right Expr
}

// Unary represents expressions with one expr and one operator, like !wololo, -12
type Unary struct {
	Token Token
	Right Expr
}

// Variable represents a variable reference in code
type Variable struct {
	Name Token
}
