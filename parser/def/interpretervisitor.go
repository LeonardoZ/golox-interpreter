package def

// ExpressionVisitor Interface
type ExpressionVisitor interface {
	VisitBinaryExpr(binary *Binary) (interface{}, *RuntimeError)
	VisitUnaryExpr(unary *Unary) (interface{}, *RuntimeError)
	VisitGroupingExpr(grouping *Grouping) (interface{}, *RuntimeError)
	VisitLiteralExpr(literal *Literal) (interface{}, *RuntimeError)
	VisitVariableExpr(variable *Variable) (interface{}, *RuntimeError)
	VisitAssignExpr(assign *Assign) (interface{}, *RuntimeError)
	VisitLogicalExpr(logical *Logical) (interface{}, *RuntimeError)
}

// StatementVisitor Interface
type StatementVisitor interface {
	VisitExpressionStmt(exprStmt *ExprStmt) *RuntimeError
	VisitPrintStmt(print *Print) *RuntimeError
	VisitVar(varStmt *Var) *RuntimeError
	VisitBlock(block *Block) *RuntimeError
	VisitIf(ifStmt *If) *RuntimeError
}

/*Expression and Statement Accepts */

// Accept def for type
func (print *Print) Accept(v StatementVisitor) *RuntimeError {
	return v.VisitPrintStmt(print)
}

// Accept def for type
func (varStmt *Var) Accept(v StatementVisitor) *RuntimeError {
	return v.VisitVar(varStmt)
}

// Accept def for type
func (exprStmt *ExprStmt) Accept(v StatementVisitor) *RuntimeError {
	return v.VisitExpressionStmt(exprStmt)
}

// Accept def for type
func (ifStmt *If) Accept(v StatementVisitor) *RuntimeError {
	return v.VisitIf(ifStmt)
}

// Accept def for type
func (block *Block) Accept(v StatementVisitor) *RuntimeError {
	return v.VisitBlock(block)
}

// Accept def for type
func (empty *EmptyExpr) Accept(v ExpressionVisitor) (interface{}, *RuntimeError) {
	return "", nil
}

// Accept def for type
func (literal *Literal) Accept(v ExpressionVisitor) (interface{}, *RuntimeError) {
	return v.VisitLiteralExpr(literal)
}

// Accept def for type
func (logical *Logical) Accept(v ExpressionVisitor) (interface{}, *RuntimeError) {
	return v.VisitLogicalExpr(logical)
}

// Accept def for type
func (grouping *Grouping) Accept(v ExpressionVisitor) (interface{}, *RuntimeError) {
	return v.VisitGroupingExpr(grouping)
}

// Accept def for type
func (binary *Binary) Accept(v ExpressionVisitor) (interface{}, *RuntimeError) {
	return v.VisitBinaryExpr(binary)
}

// Accept def for type
func (unary *Unary) Accept(v ExpressionVisitor) (interface{}, *RuntimeError) {
	return v.VisitUnaryExpr(unary)
}

// Accept def for type
func (variable *Variable) Accept(v ExpressionVisitor) (interface{}, *RuntimeError) {
	return v.VisitVariableExpr(variable)
}

// Accept def for type
func (assign *Assign) Accept(v ExpressionVisitor) (interface{}, *RuntimeError) {
	return v.VisitAssignExpr(assign)
}
