package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"loxlang/parser"
	"os"
	"strings"
)

var hadError = false

func main() {
	// (* (- 123) (group 45.67))
	expr := parser.Binary{
		&parser.Unary{
			parser.Token{
				parser.MINUS, "-", nil, 1,
			},
			&parser.Literal{
				123,
			},
		},
		parser.Token{
			parser.STAR, "*", nil, 1,
		},
		&parser.Grouping{
			&parser.Literal{
				45.67,
			},
		},
	}
	var printer parser.AstPrinter = parser.AstPrinter{}
	printer.Print(&expr)

	/*
		args := os.Args[0:]
		fmt.Println()
		if len(args) > 2 {
			fmt.Println("Usage: glox [script]")
		} else if len(args) == 2 {
			runFile(args[1])
		} else {
			runPrompt()
		} */
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
	tokens := parser.ScanTokens(content)

	for _, token := range tokens {
		fmt.Println(token)
	}
}
