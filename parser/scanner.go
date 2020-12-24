package parser

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
}

func ScanTokens(content string) []*Token {
	var tokens []*Token = []*Token{}

	return tokens
}
