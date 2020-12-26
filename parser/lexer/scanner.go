package lexer

import (
	"loxlang/parser"
	"loxlang/parser/def"
	"strconv"
)

var start, current, line int
var source []rune
var tokens []def.Token

// ScanTokens is the main function of the lexer/scanner
func ScanTokens(input string) []def.Token {
	tokens = []def.Token{}
	source = make([]rune, len(input))
	for i, r := range input {
		source[i] = r
	}
	start, current, line = 0, 0, 1

	for !isAtEnd() {
		start = current
		scanToken()
	}

	tokens = append(tokens, def.Token{def.EOF, "", nil, line})
	return tokens
}

func scanToken() {
	char := advance()

	switch char {
	case '(':
		addToken(def.LEFTPAREN)
		break
	case ')':
		addToken(def.RIGHTPAREN)
		break
	case '{':
		addToken(def.LEFTBRACE)
		break
	case '}':
		addToken(def.RIGHTBRACE)
		break
	case ',':
		addToken(def.COMMA)
		break
	case '.':
		addToken(def.DOT)
		break
	case '-':
		addToken(def.MINUS)
		break
	case '+':
		addToken(def.PLUS)
		break
	case ';':
		addToken(def.SEMICOLON)
		break
	case '*':
		addToken(def.STAR)
		break
	case '!':
		addToken(composeLexeme('=', def.BANGEQUAL, def.BANG))
		break
	case '=':
		addToken(composeLexeme('=', def.EQUALEQUAL, def.EQUAL))
		break
	case '<':
		addToken(composeLexeme('=', def.LESSEQUAL, def.LESS))
		break
	case '>':
		addToken(composeLexeme('=', def.GREATEREQUAL, def.GREATER))
		break
	case '/':
		if match('/') {
			for peek() != '\n' && !isAtEnd() {
				advance()
			}
		} else {
			addToken(def.SLASH)
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
			parser.LogError(line, "Unexpected character "+string(char))
		}
		break
	}
}

func processString() {
	for peek() != '"' && !isAtEnd() {
		if peek() == '\n' {
			line++
		}
		advance()
	}

	if isAtEnd() {
		parser.LogError(line, "Unterminated string")
		return
	}

	advance()
	realString := string(source[start+1 : current-1])
	addTokenWithLiteral(def.STRING, realString)
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

	numberAsStr := string(source[start:current])
	res, err := strconv.ParseFloat(numberAsStr, 64)
	if err != nil {
		parser.LogError(line, "Invalid number "+numberAsStr)
	}
	addTokenWithLiteral(def.NUMBER, res)
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
	text := string(source[start:current])
	tokenType, ok := def.Keywords[text]

	if !ok {
		tokenType = def.IDENTIFIER
	}

	addToken(tokenType)
}

func isAtEnd() bool {
	return current >= len(source)
}

func advance() rune {
	current++
	return source[current-1]
}

func addToken(tokenType def.TokenType) {
	addTokenWithLiteral(tokenType, nil)
}

func addTokenWithLiteral(tokenType def.TokenType, literal interface{}) {
	content := source[start:current]
	tokens = append(tokens, def.Token{tokenType, string(content), literal, line})
}

func composeLexeme(char rune, matches def.TokenType, replacement def.TokenType) def.TokenType {
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

func peekNext() rune {
	if current+1 >= len(source) {
		return '\x00'
	}
	return source[current+1]
}
