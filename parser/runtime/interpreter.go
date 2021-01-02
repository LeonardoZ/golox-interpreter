package runtime

import (
	"errors"
	"fmt"
	"loxlang/parser/def"
	"strings"
)

var env *Environment = GlobalEnvironment()

// Interpreter - implements Visitor Pattern
type Interpreter struct {
}

// Interpret Main method of Interpreter
func (i *Interpreter) Interpret(stmts []def.Stmt) {
	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			def.ReportRuntimeError(err)
			break
		}
	}
}

func (i *Interpreter) execute(stmt def.Stmt) *def.RuntimeError {
	err := stmt.Accept(i)
	return err
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

// VisitExpressionStmt Handles ExprStmt
func (i *Interpreter) VisitExpressionStmt(exprStmt *def.ExprStmt) *def.RuntimeError {
	_, err := i.evaluate(exprStmt.Expr)
	if err != nil {
		return err
	}
	return nil
}

// VisitPrintStmt Handles Print
func (i *Interpreter) VisitPrintStmt(print *def.Print) *def.RuntimeError {
	value, err := i.evaluate(print.Expr)
	if err != nil {
		return err
	}
	fmt.Println(i.stringfy(value))
	return nil
}

// VisitVar Handles Var
func (i *Interpreter) VisitVar(varStmt *def.Var) *def.RuntimeError {
	var value interface{}
	var err *def.RuntimeError
	if varStmt.Initializer != nil {
		value, err = i.evaluate(varStmt.Initializer)
		if err != nil {
			return err
		}
	}
	env.Define(varStmt.Name.Lexeme, value)
	return nil
}

// VisitVariableExpr Handles ExprStmt
func (i *Interpreter) VisitVariableExpr(variable *def.Variable) (interface{}, *def.RuntimeError) {
	return env.Get(variable.Name)
}

// VisitBlock Handles ExprStmt
func (i *Interpreter) VisitBlock(block *def.Block) *def.RuntimeError {
	i.executeBlock(block.Stmts, NewEnvironment(env))
	return nil
}

func (i *Interpreter) executeBlock(stmts []def.Stmt, outerEnv *Environment) {
	previous := env
	// goes back to the previous value
	defer func() { env = previous }()
	for _, s := range stmts {
		i.execute(s)
	}
	env = outerEnv
}

// VisitAssignExpr Handles Grouping
func (i *Interpreter) VisitAssignExpr(assign *def.Assign) (interface{}, *def.RuntimeError) {
	value, err := i.evaluate(assign)
	if err != nil {
		return nil, err
	}
	env.Assign(assign.Name, value)
	return value, nil
}

// VisitLiteralExpr Handles Literal
func (i *Interpreter) VisitLiteralExpr(literal *def.Literal) (interface{}, *def.RuntimeError) {
	return literal.Value, nil
}

// VisitGroupingExpr Handles Grouping
func (i *Interpreter) VisitGroupingExpr(grouping *def.Grouping) (interface{}, *def.RuntimeError) {
	return i.evaluate(grouping.Expression)
}

// VisitBinaryExpr Handles Binary
func (i *Interpreter) VisitBinaryExpr(binary *def.Binary) (interface{}, *def.RuntimeError) {
	var ok *def.RuntimeError
	var leftVal, rightVal float64

	right, rOk := i.evaluate(binary.Right)
	if rOk != nil {
		return nil, rOk
	}
	left, lOk := i.evaluate(binary.Left)
	if lOk != nil {
		return nil, lOk
	}

	switch binary.Token.Type {
	case def.GREATER:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)
		if ok != nil {
			return nil, ok
		}
		return leftVal > rightVal, nil
	case def.GREATEREQUAL:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			return nil, ok
		}
		return leftVal >= rightVal, nil
	case def.LESS:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			return nil, ok
		}
		return leftVal < rightVal, nil
	case def.LESSEQUAL:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			return nil, ok
		}
		return leftVal <= rightVal, nil
	case def.BANGEQUAL:
		return !i.isEqual(left, right), nil
	case def.EQUALEQUAL:
		return i.isEqual(left, right), nil
	case def.MINUS:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			return nil, ok
		}
		return leftVal - rightVal, nil
	case def.SLASH:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			return nil, ok
		}
		if rightVal == 0.0 {
			divideByZeroError := &def.RuntimeError{
				Token:   binary.Token,
				Message: "Can't divide by 0",
			}
			return nil, divideByZeroError
		}
		return leftVal / rightVal, nil
	case def.STAR:
		leftVal, rightVal, ok = i.checkNumberOperands(binary.Token, left, right)

		if ok != nil {
			return nil, ok
		}
		return leftVal * rightVal, nil
	case def.PLUS:
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
		return nil, &def.RuntimeError{
			Token:   binary.Token,
			Message: "Invalid values in + operator",
		}
	}
	return nil, nil
}

// VisitUnaryExpr Handles Unary
func (i *Interpreter) VisitUnaryExpr(unary *def.Unary) (interface{}, *def.RuntimeError) {
	right, ok := i.evaluate(unary.Right)
	if ok != nil {
		return nil, ok
	}
	switch unary.Token.Type {
	case def.BANG:
		res, compErr := i.isTruthy(right)
		if compErr != nil {
			return nil, &def.RuntimeError{
				Token:   unary.Token,
				Message: compErr.Error(),
			}
		}
		return !res, nil
	case def.MINUS:
		value, mOk := i.checkNumberOperand(unary.Token, right)
		if mOk != nil {
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

func (i Interpreter) checkNumberOperand(token def.Token, operand interface{}) (float64, *def.RuntimeError) {
	value, ok := operand.(float64)
	if !ok {
		return 0, &def.RuntimeError{
			Token:   token,
			Message: "Operand must be a number",
		}
	}
	return value, nil
}

func (i Interpreter) checkNumberOperands(token def.Token, left interface{}, right interface{}) (float64, float64, *def.RuntimeError) {
	leftVal, lOk := left.(float64)
	rightVal, rOk := right.(float64)

	if !lOk || !rOk {
		return 0, 0, &def.RuntimeError{
			Token:   token,
			Message: "Operand must be a number",
		}
	}
	return leftVal, rightVal, nil
}

func (i *Interpreter) evaluate(expr def.Expr) (interface{}, *def.RuntimeError) {
	return expr.Accept(i)
}
