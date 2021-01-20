package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"loxlang/parser"
	"loxlang/parser/def"
	"loxlang/parser/lexer"
	"loxlang/parser/pass"
	"loxlang/parser/runtime"
	"os"
	"strings"
)

func main() {
	args := os.Args[0:]
	fmt.Println()
	if len(args) > 2 {
		fmt.Println("Usage: lox [script]")
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}

func runFile(filePath string) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	run(string(dat))
	if def.HadError {
		os.Exit(65)
	}
	if def.HadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	interpreter := runtime.Interpreter{}
	fmt.Println("## GoLox REPL ##")
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.Replace(line, "\n", "", -1)
		if line == "" {
			break
		}
		tokens := lexer.ScanTokens(line)
		if def.HadError {
			def.HadError = false
			continue
		}

		stmts, _ := parser.Parse(tokens)
		if def.HadError {
			def.HadError = false
			continue
		}
		interpreter.Interpret(stmts)
		def.HadError = false
	}
}

func run(content string) {
	// lexer
	tokens := lexer.ScanTokens(content)

	if def.HadError {
		def.HadError = false
		return
	}

	// parser
	stmts, err := parser.Parse(tokens)
	if err != nil {
		def.HadError = false
		return
	}

	interpreter := runtime.NewInterpreter()

	// static analyses
	resolver := pass.NewResolver(*interpreter)
	resolver.ResolveStmts(stmts)

	if def.HadError {
		def.HadError = false
		return
	}
	// runtime
	interpreter.Interpret(stmts)

}
