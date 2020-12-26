package parser

import (
	"fmt"
	"strconv"
)

var start, current, line int
var source []rune
var tokens []Token

func ScanTokens(input string) []Token {
	tokens = []Token{}
	source = make([]rune, len(input))
	for i, r := range input {
		source[i] = r
	}
	start, current, line = 0, 0, 1

	for !isAtEnd() {
		start = current
		scanToken()
	}

	tokens = append(tokens, Token{EOF, "", nil, line})
	return tokens
}

func createRuneMap(input string) {

}

func isAtEnd() bool {
	return current >= len(source)
}

func scanToken() {
	char := advance()
	switch char {
	case '(':
		addToken(LEFTPAREN)
		break
	case ')':
		addToken(RIGHTPAREN)
		break
	case '{':
		addToken(LEFTBRACE)
		break
	case '}':
		addToken(RIGHTBRACE)
		break
	case ',':
		addToken(COMMA)
		break
	case '.':
		addToken(DOT)
		break
	case '-':
		addToken(MINUS)
		break
	case '+':
		addToken(PLUS)
		break
	case ';':
		addToken(SEMICOLON)
		break
	case '*':
		addToken(STAR)
		break
	case '!':
		addToken(composeLexeme('=', BANGEQUAL, BANG))
		break
	case '=':
		addToken(composeLexeme('=', EQUALEQUAL, EQUAL))
		break
	case '<':
		addToken(composeLexeme('=', LESSEQUAL, LESS))
		break
	case '>':
		addToken(composeLexeme('=', GREATEREQUAL, GREATER))
		break
	case '/':
		if match('/') {
			for peek() != '\n' && !isAtEnd() {
				advance()
			}
		} else {
			addToken(SLASH)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break
	case '\n':
		line++
		break

	case '"':
		processString()
		break
	default:
		if isDigit(char) {
			processNumber()
		} else if isAlpha(char) {
			processIdentifier()
		} else {
			LogError(line, "Unexpected character "+string(char))
		}
		break
	}
}

func advance() rune {
	current++
	return source[current-1]
}

func addToken(Type TokenType) {
	addTokenWithLiteral(Type, nil)
}

func addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	content := source[start:current]
	tokens = append(tokens, Token{tokenType, string(content), literal, line})
}

func composeLexeme(char rune, matches TokenType, replacement TokenType) TokenType {
	if match(char) {
		return matches
	} else {
		return replacement
	}
}

func match(expected rune) bool {
	if isAtEnd() {
		return false
	}
	currentChar := source[current]
	if currentChar != expected {
		return false
	}

	current++
	return true
}

func peek() rune {
	if isAtEnd() {
		return '\x00'
	}
	return source[current]
}

func processString() {
	for peek() != '"' && !isAtEnd() {
		if peek() == '\n' {
			line++
		}
		advance()
	}

	if isAtEnd() {
		LogError(line, "Unterminated string")
		return
	}

	advance()
	fmt.Println(start + 1)
	fmt.Println(current - 1)
	realString := string(source[start+1 : current-1])
	addTokenWithLiteral(STRING, realString)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func processNumber() {
	for isDigit(peek()) {
		advance()
	}

	if (peek() == '.') && isDigit(peekNext()) {
		advance()

		for isDigit(peek()) {
			advance()
		}
	}

	numberAsStr := string(source[start : current-1])
	res, err := strconv.ParseFloat(numberAsStr, 64)
	if err != nil {
		LogError(line, "Invalid number "+numberAsStr)
		panic(err)
	}
	addTokenWithLiteral(NUMBER, res)
}

func peekNext() rune {
	if current+1 >= len(source) {
		return '\x00'
	}
	return source[current+1]
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphanumeric(c rune) bool {
	return isAlpha(c) || isDigit((c))
}

func processIdentifier() {
	for isAlphanumeric(peek()) {
		advance()
	}
	text := string(source[start : current-1])
	tokenType, ok := Keywords[text]

	if !ok {
		tokenType = IDENTIFIER
	}

	addToken(tokenType)
}
