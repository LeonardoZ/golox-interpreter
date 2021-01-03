package runtime

import (
	"fmt"
	"loxlang/parser/def"
)

// Environment Type for handling variable name x value binding in runtime
type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

// GlobalEnvironment Creates new environment
func GlobalEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

// NewEnvironment Creates new environment
func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]interface{}),
		enclosing: enclosing,
	}
}

// Assign Assignes value to a variable
func (env *Environment) Assign(name def.Token, newValue interface{}) *def.RuntimeError {
	_, present := env.values[name.Lexeme]
	if present {
		env.values[name.Lexeme] = newValue
		return nil
	}
	if env.enclosing != nil {
		return env.enclosing.Assign(name, newValue)
	}
	return &def.RuntimeError{
		Token:   name,
		Message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
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
	if env.enclosing != nil {
		return env.enclosing.Get(token)
	}
	return nil, &def.RuntimeError{
		Token:   token,
		Message: fmt.Sprintf("Undefined variable %s", token.Lexeme),
	}
}
