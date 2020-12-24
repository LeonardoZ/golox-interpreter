package parser

import "fmt"

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line=%b] Error %s: %s\n", line, where, message)
}
