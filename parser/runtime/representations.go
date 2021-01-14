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
	FunDecl def.Function
	Closure *Environment
}

// String counts how many parameters there are in a function
func (f *CallableFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.FunDecl.Name.Lexeme)
}

// Arity counts how many parameters there are in a function
func (f *CallableFunction) Arity() int {
	return len(f.FunDecl.Params)
}

// Call invoked the function
func (f *CallableFunction) Call(i *Interpreter, args []interface{}) (interface{}, *def.RuntimeError) {
	localEnv := NewEnvironment(f.Closure)
	for i, p := range f.FunDecl.Params {
		(*localEnv).Define(p.Lexeme, args[i])
	}
	err := i.executeBlock(f.FunDecl.Body, localEnv)
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
