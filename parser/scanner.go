package parser

import "strings"

type Token struct {
	Char string
	Pos  int
}

func ScanTokens(content string) []*Token {
	var tokens []*Token
	for pos, value := range content {
		result := strings.ReplaceAll(string(value), "\n", "\\n")

		token := Token{Char: result, Pos: pos}
		tokens = append(tokens, &token)
	}
	return tokens
}
