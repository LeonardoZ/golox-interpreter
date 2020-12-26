package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"loxlang/parser"
	"loxlang/parser/def"
	"loxlang/parser/lexer"
	"os"
	"strings"
)

var hadError = false

func main() {
	args := os.Args[0:]
	fmt.Println()
	if len(args) > 2 {
		fmt.Println("Usage: glox [script]")
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
	if hadError {
		os.Exit(65)
	}

}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("> ")
		line, _ := reader.ReadString('\n')
		line = strings.Replace(line, "\n", "", -1)
		if line == "" {
			break
		}
		run(line)
		hadError = false

	}
}
func run(content string) {
	tokens := lexer.ScanTokens(content)
	result := parser.Parse(tokens)

	if hadError {
		return
	}
	ast := def.AstPrinter{}
	ast.Print(result)
}
