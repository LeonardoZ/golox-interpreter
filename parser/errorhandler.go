package parser

import (
	"errors"
	"fmt"
	"loxlang/parser/def"
)

// LogError - Logs error
func LogError(line int, message string) {
	Report(line, "", message)
}

// Report - Log error with more info
func Report(line int, where string, message string) {
	fmt.Printf("[line=%b] Error %s: %s\n", line, where, message)
}

// Generates an Error
func CreateError(token def.Token, message string) error {
	if token.Type == def.EOF {
		Report(token.Line, " at end ", message)
	} else {
		Report(token.Line, " at '"+token.Lexeme+"'", message)
	}
	return errors.New(message)
}
