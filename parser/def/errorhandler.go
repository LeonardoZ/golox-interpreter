package def

import (
	"errors"
	"fmt"
)

var HadError bool = false
var HadRuntimeError bool = false

// LogError - Logs error
func LogError(line int, message string) {
	Report(line, "", message)
}

// Report - Log error with more info
func Report(line int, where string, message string) {
	fmt.Printf("[line=%b] Error %s: %s\n", line, where, message)
}

func ReportRuntimeError(runtimeError *RuntimeError) {
	HadRuntimeError = true
	fmt.Printf(runtimeError.Error())
}

// Generates an Error
func CreateError(token Token, message string) error {
	if token.Type == EOF {
		Report(token.Line, " at end ", message)
	} else {
		Report(token.Line, " at '"+token.Lexeme+"'", message)
	}
	return errors.New(message)
}
