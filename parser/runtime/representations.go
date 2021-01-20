package runtime

import (
	"fmt"
	"loxlang/parser/def"
)

// Callable is a function-call representation in runtime
type Callable interface {
	Arity() int
	Call(i *Interpreter, args []interface{}) (interface{}, *def.RuntimeError)
}

// CallableFunction is a concrete representation of a user-defined function to be called
type CallableFunction struct {
	Name         string
	FunctionExpr def.FunctionExpr
	Closure      *Environment
}

// String counts how many parameters there are in a function
func (f *CallableFunction) String() string {
	if f.Name != "" {
		return fmt.Sprintf("<fn %s>", f.Name)
	}
	return "<fn>"
}

// Arity counts how many parameters there are in a function
func (f *CallableFunction) Arity() int {
	return len(f.FunctionExpr.Params)
}

// Call invoked the function
func (f *CallableFunction) Call(i *Interpreter, args []interface{}) (interface{}, *def.RuntimeError) {
	localEnv := NewEnvironment(f.Closure)
	for i := range f.FunctionExpr.Params {
		(*localEnv).Define(args[i])
	}
	err := i.executeBlock(f.FunctionExpr.Body, localEnv)
	if err != nil {
		if err.Type == def.RETURNSTMT {
			return err.Value, nil
		}
		return nil, err
	}
	return nil, nil
}

// ReturnValue represents the value that returns from a function
type ReturnValue struct {
	Value interface{}
}
