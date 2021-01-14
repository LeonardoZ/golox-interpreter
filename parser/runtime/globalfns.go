package runtime

import (
	"loxlang/parser/def"
	"time"
)

// ClockCallable default clock function implementation
type ClockCallable struct{}

// Arity of the clock fn
func (c *ClockCallable) Arity() int {
	return 0
}

// Call representation of the clock fn
func (c *ClockCallable) Call(i *Interpreter, args []interface{}) (interface{}, *def.RuntimeError) {
	return time.Now().UnixNano() / int64(time.Millisecond), nil
}
