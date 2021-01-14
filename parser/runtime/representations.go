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
	localEnv := NewEnvironment(Globals)
	for i, p := range f.FunDecl.Params {
		(*localEnv).Define(p.Lexeme, args[i])
	}
	err := i.executeBlock(f.FunDecl.Body, localEnv)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
