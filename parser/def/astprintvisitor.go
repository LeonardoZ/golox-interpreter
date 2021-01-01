package def

// StrVisitor Interface
type StrVisitor interface {
	VisitBinaryExprStr(binary *Binary) string
	VisitUnaryExpStr(unary *Unary) string
	VisitGroupingExprStr(grouping *Grouping) string
	VisitLiteralExprStr(literal *Literal) string
	VisitVariableExprStr(variable *Variable) string
}

// AcceptStr def for type
func (empty *EmptyExpr) AcceptStr(v StrVisitor) string {
	return ""
}

// AcceptStr def for type
func (literal *Literal) AcceptStr(v StrVisitor) string {
	return v.VisitLiteralExprStr(literal)
}

// AcceptStr def for type
func (grouping *Grouping) AcceptStr(v StrVisitor) string {
	return v.VisitGroupingExprStr(grouping)
}

// AcceptStr def for type
func (binary *Binary) AcceptStr(v StrVisitor) string {
	return v.VisitBinaryExprStr(binary)
}

// AcceptStr def for type
func (unary *Unary) AcceptStr(v StrVisitor) string {
	return v.VisitUnaryExpStr(unary)
}

// AcceptStr def for type
func (variable *Variable) AcceptStr(v StrVisitor) string {
	return v.VisitVariableExprStr(variable)
}