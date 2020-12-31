package def

import (
	"errors"
	"fmt"
	"strings"
)

type InterpreterVisitor interface {
	visitBinaryExpr(binary *Binary) (interface{}, *RuntimeError)
	visitUnaryExpr(unary *Unary) (interface{}, *RuntimeError)
	visitGroupingExpr(grouping *Grouping) (interface{}, *RuntimeError)
	visitLiteralExpr(literal *Literal) (interface{}, *RuntimeError)
}

// Interpreter - implements Visitor Pattern
type Interpreter struct {
}

// Interpret Main method of Interpreter
func (i *Interpreter) Interpret(expr Expr) string {
	value, err := i.Evaluate(expr)
	if err != nil {
		ReportRuntimeError(err)
		return ""
	}
	return i.stringfy(value)

}
func (i *Interpreter) stringfy(value interface{}) string {
	if value == nil {
		return ""
	}
	parsed, isFloat := value.(float64)
	if isFloat {
		res := fmt.Sprintf("%f", parsed)
		if strings.HasSuffix(res, ".000000") {
			text := res[0 : len(res)-7]
			return text
		}
		return res
	}

	return fmt.Sprintf("%v", value)
}

func (i *Interpreter) visitLiteralExpr(literal *Literal) (interface{}, *RuntimeError) {
	return literal.Value, nil
}

func (i *Interpreter) visitGroupingExpr(grouping *Grouping) (interface{}, *RuntimeError) {
	return i.Evaluate(grouping.Expression)
}

func (i *Interpreter) visitBinaryExpr(binary *Binary) (interface{}, *RuntimeError) {
	var ok *RuntimeError
	var leftVal, rightVal float64

	right, rOk := i.Evaluate(binary.Right)
	if rOk != nil {
		ReportRuntimeError(rOk)
		return nil, rOk
	}
	left, lOk := i.Evaluate(binary.Left)
	if lOk != nil {
		ReportRuntimeError(lOk)
		return nil, lOk
	}

	switch binary.Token.Type {
	case GREATER:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)
		if ok != nil {
			ReportRuntimeError(ok)
			return nil, ok
		}
		return leftVal > rightVal, nil
	case GREATEREQUAL:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			ReportRuntimeError(ok)
			return nil, ok
		}
		return leftVal >= rightVal, nil
	case LESS:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			ReportRuntimeError(ok)
			return nil, ok
		}
		return leftVal < rightVal, nil
	case LESSEQUAL:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			ReportRuntimeError(ok)
			return nil, ok
		}
		return leftVal <= rightVal, nil
	case BANGEQUAL:
		return !i.isEqual(left, right), nil
	case EQUALEQUAL:
		return i.isEqual(left, right), nil
	case MINUS:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			ReportRuntimeError(ok)
			return nil, ok
		}
		return leftVal - rightVal, nil
	case SLASH:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			ReportRuntimeError(ok)
			return nil, ok
		}
		return leftVal / rightVal, nil
	case STAR:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			ReportRuntimeError(ok)
			return nil, ok
		}
		return leftVal * rightVal, nil
	case PLUS:
		floatLeft, isFloatLeft := left.(float64)
		floatRight, isFloatRight := right.(float64)
		if isFloatLeft && isFloatRight {
			return floatLeft + floatRight, nil
		}

		stringRight, isStringRight := right.(string)
		stringLeft, isStringLeft := left.(string)
		if isStringLeft && isStringRight {
			return stringLeft + stringRight, nil
		}
		return nil, &RuntimeError{
			Token:   binary.Token,
			Message: "Invalid values in + operator",
		}
	}
	return nil, nil
}

func (i *Interpreter) visitUnaryExpr(unary *Unary) (interface{}, *RuntimeError) {
	right, ok := i.Evaluate(unary.Right)
	if ok != nil {
		ReportRuntimeError(ok)
		return nil, ok
	}
	switch unary.Token.Type {
	case BANG:
		res, compErr := i.isTruthy(right)
		if compErr != nil {
			ReportRuntimeError(&RuntimeError{
				Token:   unary.Token,
				Message: compErr.Error(),
			})
		}
		return !res, nil
	case MINUS:
		value, mOk := i.checkNumberOperand(unary.Token, right)
		if mOk != nil {
			ReportRuntimeError(mOk)
			return nil, mOk
		}
		return -value, nil
	}
	return nil, nil
}

func (i Interpreter) isTruthy(value interface{}) (bool, error) {
	if value == nil {
		return false, nil
	}
	if boolval, ok := value.(bool); ok {
		return boolval, nil
	}
	return false, errors.New("Not a boolean")
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

func (i Interpreter) checkNumberOperand(token Token, operand interface{}) (float64, *RuntimeError) {
	value, ok := operand.(float64)
	if !ok {
		return 0, &RuntimeError{
			Token:   token,
			Message: "Operand must be a number",
		}
	}
	return value, nil
}

func (i Interpreter) checkNumberOperands(token Token, left interface{}, right interface{}) (float64, float64, *RuntimeError) {
	leftVal, lOk := left.(float64)
	rightVal, rOk := right.(float64)

	if !lOk || !rOk {
		return 0, 0, &RuntimeError{
			Token:   token,
			Message: "Operand must be a number",
		}
	}
	return leftVal, rightVal, nil
}

func (i *Interpreter) Evaluate(expr Expr) (interface{}, *RuntimeError) {
	return expr.Accept(i)
}

func (empty *EmptyExpr) Accept(v InterpreterVisitor) (interface{}, *RuntimeError) {
	return "", nil
}

func (literal *Literal) Accept(v InterpreterVisitor) (interface{}, *RuntimeError) {
	return v.visitLiteralExpr(literal)
}

func (grouping *Grouping) Accept(v InterpreterVisitor) (interface{}, *RuntimeError) {
	return v.visitGroupingExpr(grouping)
}

func (binary *Binary) Accept(v InterpreterVisitor) (interface{}, *RuntimeError) {
	return v.visitBinaryExpr(binary)
}

func (unary *Unary) Accept(v InterpreterVisitor) (interface{}, *RuntimeError) {
	return v.visitUnaryExpr(unary)
}
