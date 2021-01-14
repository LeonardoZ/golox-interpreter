package def

// Stmt Statement
type Stmt interface {
	// Accept Method for StatementVisitor
	Accept(visitor StatementVisitor) *RuntimeError
}

// ExprStmt Expression Statements
type ExprStmt struct {
	Expr Expr
}

// Block Expression Statements
type Block struct {
	Stmts []Stmt
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

// If represents conditional if statements
type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

// Function represents a function declaration
type Function struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

// While represents repetition loop
type While struct {
	Condition Expr
	Body      Stmt
}

// ControlFlow represents break or continue
type ControlFlow struct {
	Type ErrorType
}

// Expr Mostly generic Tree Node
type Expr interface {
	AcceptStr(visitor StrVisitor) string
	Accept(interpreterVisitor ExpressionVisitor) (interface{}, *RuntimeError)
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

// Call represents a function call
type Call struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

// Binary represents expressions with two expr and one operator, like 1 + 2, a > b
type Binary struct {
	Left  Expr
	Token Token
	Right Expr
}

// Logical represents 'and' and 'or' expressions
type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
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

// Assign represents variable assign
type Assign struct {
	Name  Token
	Value Expr
}
