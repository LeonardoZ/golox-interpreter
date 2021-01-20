package pass

import (
	"fmt"
	"loxlang/parser/def"
	"loxlang/parser/runtime"
)

type fnScope int

// Identifies if it is in function scope or not
const (
	ScopeNone fnScope = iota
	ScopeFunction
)

// Resolver Walks the parse tree doing static analyses - variable resolution
type Resolver struct {
	Interpreter runtime.Interpreter
	Scopes      ScopeStack
	CurrentSope fnScope
}

// NewResolver creates new instance of resolver
func NewResolver(i runtime.Interpreter) (r *Resolver) {
	return &Resolver{
		Interpreter: i,
		Scopes:      ScopeStack{},
		CurrentSope: ScopeNone,
	}
}

// VisitBlock Handles ExprStmt
func (r *Resolver) VisitBlock(block *def.Block) *def.RuntimeError {
	r.beginScope()
	r.ResolveStmts(block.Stmts)
	r.endScope()
	return nil
}

// VisitVar Handles Var
func (r *Resolver) VisitVar(varStmt *def.Var) *def.RuntimeError {
	r.declare(varStmt.Name)
	if varStmt.Initializer != nil {
		err := r.resolveExpr(varStmt.Initializer)
		if err != nil {
			return err
		}
	}
	r.define(varStmt.Name)
	return nil
}

// VisitVariableExpr Handles ExprStmt
func (r *Resolver) VisitVariableExpr(variable *def.Variable) (interface{}, *def.RuntimeError) {
	notIsEmpty := !r.Scopes.IsEmpty()
	scope, _ := r.Scopes.Peek()
	scopeValue, containsKey := scope[variable.Name.Lexeme]

	if notIsEmpty && containsKey && !scopeValue.IsDefined {
		return nil, &def.RuntimeError{
			Token:   variable.Name,
			Message: fmt.Sprintf("%s: Can't read local variable in its own initializer.", variable.Name.Lexeme),
		}
	}
	r.resolveLocal(variable, variable.Name)
	return nil, nil
}

// VisitAssignExpr Handles Grouping
func (r *Resolver) VisitAssignExpr(assign *def.Assign) (interface{}, *def.RuntimeError) {
	err := r.resolveExpr(assign.Value)
	if err != nil {
		return nil, err
	}
	r.resolveLocal(assign, assign.Name)

	return nil, nil
}

// VisitFunction Handles Function
func (r *Resolver) VisitFunction(function *def.Function) *def.RuntimeError {
	r.declare(function.Name)
	r.define(function.Name)
	r.resolveFunction(*&function.FuncExpr, ScopeFunction)
	return nil
}

// VisitExpressionStmt Handles ExprStmt
func (r *Resolver) VisitExpressionStmt(exprStmt *def.ExprStmt) *def.RuntimeError {
	err := r.resolveExpr(exprStmt.Expr)
	if err != nil {
		return err
	}
	return nil
}

// VisitIf Handles Grouping
func (r *Resolver) VisitIf(ifStmt *def.If) *def.RuntimeError {
	err := r.resolveExpr(ifStmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(ifStmt.ThenBranch)
	if err != nil {
		return err
	}
	if ifStmt.ElseBranch != nil {
		err = r.resolveStmt(ifStmt.ElseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) beginScope() {
	r.Scopes.Push(map[string]*Variable{})
}

// ResolveStmts resolve all statements
func (r *Resolver) ResolveStmts(stmts []def.Stmt) {
	for _, s := range stmts {
		err := func(stmt def.Stmt) *def.RuntimeError {
			return r.resolveStmt(stmt)
		}(s)
		if err != nil {
			def.ReportRuntimeError(err)
			return
		}
	}
}

func (r *Resolver) resolveStmt(stmt def.Stmt) *def.RuntimeError {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr def.Expr) *def.RuntimeError {
	_, err := expr.Accept(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) resolveFunction(function def.FunctionExpr, scope fnScope) {
	enclosingScope := r.CurrentSope
	r.CurrentSope = scope
	r.beginScope()
	for _, p := range function.Params {
		func(param def.Token) {
			r.declare(param)
			r.define(param)
		}(p)
	}
	r.ResolveStmts(function.Body)
	r.endScope()
	r.CurrentSope = enclosingScope
}

func (r *Resolver) resolveLocal(expr def.Expr, token def.Token) {
	for i := len(r.Scopes) - 1; i >= 0; i-- {
		scope := r.Scopes[i]
		if variable, ok := scope[token.Lexeme]; ok {
			// doesn't need anon func because it stops after it finds the first one
			r.Interpreter.Resolve(expr, len(r.Scopes)-1-i, variable.Slot)
			return
		}
	}
}

func (r *Resolver) declare(token def.Token) {
	if r.Scopes.IsEmpty() {
		return
	}
	scope, _ := r.Scopes.Peek()
	if _, ok := scope[token.Lexeme]; ok {
		def.CreateError(token, "Already variable with this name in this scope.")
	}
	scope[token.Lexeme] = &Variable{
		IsDefined: false,
		Slot:      len(scope),
	}
}

func (r *Resolver) define(token def.Token) {
	if r.Scopes.IsEmpty() {
		return
	}
	scope, _ := r.Scopes.Peek()
	if variable, ok := scope[token.Lexeme]; ok {
		variable.IsDefined = true
	}
}

func (r *Resolver) endScope() {
	r.Scopes.Pop()
}

// VisitPrintStmt Handles Print
func (r *Resolver) VisitPrintStmt(print *def.Print) *def.RuntimeError {
	r.resolveExpr(print.Expr)
	return nil
}

// VisitWhile Handles Grouping
func (r *Resolver) VisitWhile(whileStmt *def.While) *def.RuntimeError {
	err := r.resolveExpr(whileStmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(whileStmt.Body)
	if err != nil {
		return err
	}
	return nil
}

// VisitControlFlow Handles Grouping
func (r *Resolver) VisitControlFlow(controlFlow *def.ControlFlow) *def.RuntimeError {
	return nil
}

// VisitFunctionExpr Handles anonymous functions
func (r *Resolver) VisitFunctionExpr(function *def.FunctionExpr) (interface{}, *def.RuntimeError) {
	r.resolveFunction(*function, ScopeFunction)
	return nil, nil
}

// VisitReturnStmt Handles Return inside function
func (r *Resolver) VisitReturnStmt(returnStmt *def.Return) *def.RuntimeError {
	if r.CurrentSope == ScopeNone {
		return &def.RuntimeError{
			Token:   returnStmt.Keyword,
			Message: "Can't return from top-level code",
		}
	}
	if returnStmt.Value != nil {
		err := r.resolveExpr(returnStmt.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// VisitLiteralExpr Handles Literal
func (r *Resolver) VisitLiteralExpr(literal *def.Literal) (interface{}, *def.RuntimeError) {
	return nil, nil
}

// VisitGroupingExpr Handles Grouping
func (r *Resolver) VisitGroupingExpr(grouping *def.Grouping) (interface{}, *def.RuntimeError) {
	r.resolveExpr(grouping.Expression)
	return nil, nil
}

// VisitBinaryExpr Handles Binary
func (r *Resolver) VisitBinaryExpr(binary *def.Binary) (interface{}, *def.RuntimeError) {
	err := r.resolveExpr(binary.Left)
	if err != nil {
		return nil, err
	}
	err = r.resolveExpr(binary.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// VisitLogicalExpr Handles Logical
func (r *Resolver) VisitLogicalExpr(logical *def.Logical) (interface{}, *def.RuntimeError) {
	err := r.resolveExpr(logical.Left)
	if err != nil {
		return nil, err
	}
	err = r.resolveExpr(logical.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// VisitUnaryExpr Handles Unary
func (r *Resolver) VisitUnaryExpr(unary *def.Unary) (interface{}, *def.RuntimeError) {
	err := r.resolveExpr(unary.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// VisitCallExpr Handles Call expressions, like ()
func (r *Resolver) VisitCallExpr(call *def.Call) (interface{}, *def.RuntimeError) {
	err := r.resolveExpr(call.Callee)
	if err != nil {
		return nil, err
	}
	for _, a := range call.Arguments {
		err := func(arg def.Expr) *def.RuntimeError {
			return r.resolveExpr(arg)
		}(a)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
