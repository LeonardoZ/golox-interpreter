package def

import "fmt"

// RuntimeError is used to represent any kind o error, like type mismatch, found in runtime
type RuntimeError struct {
	Token   Token
	Message string
}

func (err *RuntimeError) Error() string {
	return fmt.Sprintf("[line=%b] Error %s", err.Token.Line, err.Message)
}
