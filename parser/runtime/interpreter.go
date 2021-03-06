package runtime

import (
	"errors"
	"fmt"
	"loxlang/parser/def"
	"strings"
)

// Interpreter - implements Visitor Pattern
type Interpreter struct {
	Globals map[string]interface{}
	Env     *Environment
	Locals  map[def.Expr]int
	Slots   map[def.Expr]int
}

// NewInterpreter creates and sets up new Interpreter
func NewInterpreter() *Interpreter {
	globals := map[string]interface{}{}
	globals["clock"] = &ClockCallable{}
	return &Interpreter{
		Globals: globals,
		Locals:  map[def.Expr]int{},
		Slots:   map[def.Expr]int{},
	}
}

// Interpret Main method of Interpreter
func (i *Interpreter) Interpret(stmts []def.Stmt) {
	for _, s := range stmts {
		err := func(stmt def.Stmt) *def.RuntimeError {
			return i.execute(stmt)
		}(s)
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

// Resolve todo
func (i *Interpreter) Resolve(expr def.Expr, scopeDepth, slot int) *def.RuntimeError {
	i.Locals[expr] = scopeDepth
	i.Slots[expr] = slot
	return nil
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
	i.define(varStmt.Name, value)

	return nil
}

// VisitVariableExpr Handles ExprStmt
func (i *Interpreter) VisitVariableExpr(variable *def.Variable) (interface{}, *def.RuntimeError) {
	return i.lookupVariable(variable.Name, variable)
}

func (i *Interpreter) lookupVariable(name def.Token, expr def.Expr) (interface{}, *def.RuntimeError) {
	distance, ok := i.Locals[expr]
	if ok {
		return i.Env.GetAt(distance, i.Slots[expr])
	}
	if globalVal, ok := i.Globals[name.Lexeme]; ok {
		return globalVal, nil
	}
	return nil, &def.RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("Undefined variable %s.", name.Lexeme),
	}
}

// VisitBlock Handles ExprStmt
func (i *Interpreter) VisitBlock(block *def.Block) *def.RuntimeError {
	err := i.executeBlock(block.Stmts, NewEnvironment(i.Env))
	if err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) executeBlock(stmts []def.Stmt, outerEnv *Environment) *def.RuntimeError {
	previous := i.Env
	// goes back to the previous value
	defer func() { i.Env = previous }()
	i.Env = outerEnv
	for _, s := range stmts {
		err := func(stmt def.Stmt) *def.RuntimeError {
			return i.execute(s)
		}(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// VisitAssignExpr Handles AssignExpr
func (i *Interpreter) VisitAssignExpr(assign *def.Assign) (interface{}, *def.RuntimeError) {
	value, err := i.evaluate(assign.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := i.Locals[assign]
	slot, okSlot := i.Slots[assign]

	if ok && okSlot {
		i.Env.AssignAt(distance, value, slot)
		return value, nil
	}
	if _, ok = i.Globals[assign.Name.Lexeme]; ok {
		i.Globals[assign.Name.Lexeme] = value
		return value, nil
	}
	return nil, &def.RuntimeError{
		Token:   assign.Name,
		Message: fmt.Sprintf("Undefined variable %s.", assign.Name.Lexeme),
	}

}

// VisitIf Handles Grouping
func (i *Interpreter) VisitIf(ifStmt *def.If) *def.RuntimeError {
	condition, err := i.evaluate(ifStmt.Condition)
	if err != nil {
		return err
	}
	result, truthyErr := i.isTruthy(condition)
	if truthyErr != nil {
		return &def.RuntimeError{
			Token:   def.Token{},
			Message: "Error evaluating isTruthy",
		}
	}
	if result {
		err = i.execute(ifStmt.ThenBranch)
		if err != nil {
			return err
		}
	} else if ifStmt.ElseBranch != nil {
		err = i.execute(ifStmt.ElseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

// VisitWhile Handles Grouping
func (i *Interpreter) VisitWhile(whileStmt *def.While) *def.RuntimeError {
	for {
		condition, err := i.evaluate(whileStmt.Condition)
		if err != nil {
			return err
		}
		result, truthyErr := i.isTruthy(condition)
		if truthyErr != nil {
			return &def.RuntimeError{
				Token:   def.Token{},
				Message: "Error evaluating isTruthy",
			}
		}

		if !result {
			return nil
		}

		err = i.execute(whileStmt.Body)
		if err != nil {
			if err.Type == def.CONTROLFLOWBREAK {
				break
			} else {
				return err
			}
		}

	}
	return nil
}

// VisitControlFlow Handles Grouping
func (i *Interpreter) VisitControlFlow(controlFlow *def.ControlFlow) *def.RuntimeError {
	return &def.RuntimeError{
		Type: controlFlow.Type,
	}
}

// VisitFunction Handles Function
func (i *Interpreter) VisitFunction(function *def.Function) *def.RuntimeError {
	callable := CallableFunction{Name: function.Name.Lexeme, FunctionExpr: function.FuncExpr, Closure: i.Env}
	i.define(function.Name, callable)
	return nil
}

// VisitFunctionExpr Handles anonymous functions
func (i *Interpreter) VisitFunctionExpr(function *def.FunctionExpr) (interface{}, *def.RuntimeError) {
	return CallableFunction{Name: "", FunctionExpr: *function, Closure: i.Env}, nil
}

// VisitReturnStmt Handles Return inside function
func (i *Interpreter) VisitReturnStmt(returnStmt *def.Return) *def.RuntimeError {
	var value interface{}
	var err *def.RuntimeError
	if returnStmt.Value != nil {
		value, err = i.evaluate(returnStmt.Value)
		if err != nil {
			return err
		}
	}
	return &def.RuntimeError{
		Type:  def.RETURNSTMT,
		Value: value,
	}
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

// VisitLogicalExpr Handles Logical
func (i *Interpreter) VisitLogicalExpr(logical *def.Logical) (interface{}, *def.RuntimeError) {
	left, err := i.evaluate(logical.Left)
	if err != nil {
		return nil, err
	}
	result, tErr := i.isTruthy(left)
	if tErr != nil {
		return nil, &def.RuntimeError{
			Token:   logical.Operator,
			Message: "Invalid isTruthy boolean parse",
		}
	}

	if logical.Operator.Type == def.OR {
		if result {
			return left, nil
		}
	} else {
		if !result {
			return left, nil
		}
	}
	right, rErr := i.evaluate(logical.Right)
	if rErr != nil {
		return nil, rErr
	}
	return right, nil
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

// VisitCallExpr Handles Call expressions, like ()
func (i *Interpreter) VisitCallExpr(call *def.Call) (interface{}, *def.RuntimeError) {
	callee, err := i.evaluate(call.Callee)
	if err != nil {
		return nil, err
	}

	args := []interface{}{}
	for _, a := range call.Arguments {
		arg, argErr := func(aParam def.Expr) (interface{}, *def.RuntimeError) {
			return i.evaluate(aParam)
		}(a)
		if argErr != nil {
			return nil, argErr
		}
		args = append(args, arg)
	}
	callable, ok := callee.(CallableFunction)
	if !ok {
		return nil, &def.RuntimeError{
			Token:   call.Paren,
			Message: "Can only call functions and classes",
		}
	}

	if len(args) != callable.Arity() {
		return nil, &def.RuntimeError{
			Token:   call.Paren,
			Message: fmt.Sprintf("Expeted %d argumentos, but got %d", callable.Arity(), len(args)),
		}
	}
	return callable.Call(i, args)
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

func (i *Interpreter) define(name def.Token, value interface{}) {
	if i.Env != nil {
		i.Env.Define(value)
	} else {
		i.Globals[name.Lexeme] = value
	}
}
