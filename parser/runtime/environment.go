package runtime

import (
	"fmt"
	"loxlang/parser/def"
)

// Environment Type for handling variable name x value binding in runtime
type Environment struct {
	values map[string]interface{}
}

// NewEnvironment Creates new environment
func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

// Define Defines a new name x value variable
func (env *Environment) Define(name string, value interface{}) {
	env.values[name] = value
}

// Get Returns value of a name
func (env *Environment) Get(token def.Token) (interface{}, *def.RuntimeError) {
	value, present := env.values[token.Lexeme]
	if present {
		return value, nil
	}
	return nil, &def.RuntimeError{
		Token:   token,
		Message: fmt.Sprintf("Undefined variable %s", token.Lexeme),
	}
}
