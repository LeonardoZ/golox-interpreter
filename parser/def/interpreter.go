package def

type InterpreterVisitor interface {
	visitBinaryExpr(binary *Binary) string
	visitUnaryExpr(unary *Unary) string
	visitGroupingExpr(grouping *Grouping) string
	visitLiteralExpr(literal *Literal) string
}

// Interpreter - implements Visitor Pattern
type Interpreter struct {
}

func (i *Interpreter) visitLiteralExpr(literal *Literal) interface{} {
	return literal.Value
}

func (i *Interpreter) visitGroupingExpr(grouping *Grouping) interface{} {
	return i.evaluate(grouping.Expression)
}

func (i *Interpreter) visitBinaryExpr(binary *Binary) interface{} {
	right := i.evaluate(binary.Right)
	left := i.evaluate(binary.Left)

	switch binary.Token.Type {
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATEREQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESSEQUAL:
		return left.(float64) <= right.(float64)
	case BANGEQUAL:
		return !i.isEqual(left, right)
	case EQUALEQUAL:
		return i.isEqual(left, right)
	case MINUS:
		return left.(float64) - right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case PLUS:
		floatLeft, isFloatLeft := left.(float64)
		floatRight, isFloatRight := right.(float64)
		if isFloatLeft && isFloatRight {
			return floatLeft + floatRight
		}

		stringRight, isStringRight := right.(string)
		stringLeft, isStringLeft := left.(string)
		if isStringLeft && isStringRight {
			return stringLeft + stringRight
		}
		break
	}
	return nil
}

func (i *Interpreter) visitUnaryExpr(unary *Unary) interface{} {
	right := i.evaluate(unary.Right)
	switch unary.Token.Type {
	case BANG:
		return !i.isTruthy(right)
		break
	case MINUS:
		return -right.(int)
		break
	}
	return nil
}

func (i Interpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if boolval, ok := value.(bool); ok {
		return boolval
	}
	panic("Not a boolean")
}

func (i Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (empty *EmptyExpr) Accept(v InterpreterVisitor) interface{} {
	return ""
}

func (literal *Literal) Accept(v InterpreterVisitor) interface{} {
	return v.visitLiteralExpr(literal)
}

func (grouping *Grouping) Accept(v InterpreterVisitor) interface{} {
	return v.visitGroupingExpr(grouping)
}

func (binary *Binary) Accept(v InterpreterVisitor) interface{} {
	return v.visitBinaryExpr(binary)
}

func (unary *Unary) Accept(v InterpreterVisitor) interface{} {
	return v.visitUnaryExpr(unary)
}
