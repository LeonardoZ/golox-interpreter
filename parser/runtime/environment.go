package runtime

import (
	"loxlang/parser/def"
)

// Environment Type for handling variable name x value binding in runtime
type Environment struct {
	values    []interface{}
	enclosing *Environment
}

// GlobalEnvironment Creates new environment
func GlobalEnvironment() *Environment {
	return &Environment{
		values: []interface{}{},
	}
}

// NewEnvironment Creates new environment
func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    []interface{}{},
		enclosing: enclosing,
	}
}

// AssignAt Assigns to a name x the value at the variable on specified distance
func (env *Environment) AssignAt(distance int, value interface{}, slot int) {
	env.ancestor(distance).values[slot] = value
}

// Define Defines a new name x value variable
func (env *Environment) Define(value interface{}) {
	env.values = append(env.values, value)
}

// GetAt Returns value on the specified distance
func (env *Environment) GetAt(distance int, slot int) (interface{}, *def.RuntimeError) {
	vars := env.ancestor(distance).values
	found := vars[slot]
	return found, nil
}

func (env *Environment) ancestor(distance int) *Environment {
	var searchedEnv *Environment = env
	for i := 0; i < distance; i++ {
		searchedEnv = searchedEnv.enclosing
	}
	return searchedEnv
}

/*
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
*/
