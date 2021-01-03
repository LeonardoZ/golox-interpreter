package def

import "fmt"

// ErrorType Represents basic Error types
type ErrorType int8

// RuntimeError is used to represent any kind o error, like type mismatch, found in runtime
type RuntimeError struct {
	Token   Token
	Message string
	Type    ErrorType
}

func (err *RuntimeError) Error() string {
	return fmt.Sprintf("[line=%d] Error: %s", err.Token.Line, err.Message)
}

// Error Types
const (
	NORMAL           ErrorType = 0
	CONTROLFLOWBREAK ErrorType = 1
)
