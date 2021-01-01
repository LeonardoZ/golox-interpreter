package parser

import (
	"loxlang/parser/def"
)

var tokens []def.Token
var stmts []def.Stmt
var current int

// Parse is cool
func Parse(input []def.Token) []def.Stmt {
	tokens = input
	for !isAtEnd() {
		stmt, err := declaration()
		if err != nil {
			return []def.Stmt{}
		}
		stmts = append(stmts, stmt)
	}
	return stmts
}

func declaration() (def.Stmt, error) {
	if match(def.VAR) {
		varStmt, err := varDeclaration()
		if err != nil {
			synchronize()
		}
		return varStmt, nil
	}
	stmt, err := statement()
	if err != nil {
		synchronize()
	}
	return stmt, nil
}

func varDeclaration() (def.Stmt, error) {
	name, err := consume(def.IDENTIFIER, "Expect variable name")
	if err != nil {
		return nil, err
	}
	var initializer def.Expr
	if match(def.EQUAL) {
		initializer, err = expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = consume(def.SEMICOLON, "Expect ; after variable decalration")
	if err != nil {
		return nil, err
	}
	return &def.Var{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func statement() (def.Stmt, error) {
	if match(def.PRINT) {
		return printStatement()
	}
	if match(def.LEFTBRACE) {
		stmts, err := block()
		if err != nil {
			return nil, err
		}
		return &def.Block{
			Stmts: stmts,
		}, nil
	}
	return expressionStatement()
}

func block() ([]def.Stmt, error) {
	stmts := []def.Stmt{}
	for !check(def.RIGHTBRACE) && !isAtEnd() {
		stmt, err := declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	_, err := consume(def.RIGHTBRACE, "Expect '{' after block")
	if err != nil {
		return nil, err
	}
	return stmts, nil
}

func printStatement() (def.Stmt, error) {
	expr, err := expression()
	if err != nil {
		return nil, err
	}
	_, consErr := consume(def.SEMICOLON, "Expect ';' after value")
	if consErr != nil {
		return nil, consErr
	}
	return &def.Print{
		Expr: expr,
	}, nil
}

func expressionStatement() (def.Stmt, error) {
	expr, err := expression()
	if err != nil {
		return nil, err
	}
	_, consErr := consume(def.SEMICOLON, "Expect ';' after value")
	if consErr != nil {
		return nil, consErr
	}
	return &def.ExprStmt{
		Expr: expr,
	}, nil
}

func expression() (def.Expr, error) {
	return assignment()
}

func assignment() (def.Expr, error) {
	expr, err := equality()
	if err != nil {
		return nil, err
	}

	if match(def.EQUAL) {
		equals := previous()
		value, assignErr := assignment()

		if assignErr != nil {
			return nil, assignErr
		}
		if variable, res := expr.(*def.Variable); res {
			name := variable.Name
			return &def.Assign{
				Name:  name,
				Value: value,
			}, nil
		}
		reportError(equals, "Invalid assign target")
	}
	return expr, nil
}

func equality() (def.Expr, error) {
	expr, err := comparison()
	if err != nil {
		return &def.EmptyExpr{}, err
	}
	for match(def.BANGEQUAL, def.EQUALEQUAL) {
		operator := previous()
		right, errInside := comparison()
		if errInside != nil {
			return &def.EmptyExpr{}, err
		}
		expr = &def.Binary{
			Left:  expr,
			Token: operator,
			Right: right,
		}
	}
	return expr, nil
}

func comparison() (def.Expr, error) {
	expr, err := term()
	if err != nil {
		return &def.EmptyExpr{}, err
	}
	for match(def.GREATER, def.GREATEREQUAL, def.LESS, def.LESSEQUAL) {
		operator := previous()
		right, errInside := term()
		if errInside != nil {
			return &def.EmptyExpr{}, err
		}
		expr = &def.Binary{
			Left:  expr,
			Token: operator,
			Right: right,
		}
	}
	return expr, nil
}

func term() (def.Expr, error) {
	expr, err := factor()
	if err != nil {
		return &def.EmptyExpr{}, err
	}
	for match(def.MINUS, def.PLUS) {
		operator := previous()
		right, errInside := factor()
		if errInside != nil {
			return &def.EmptyExpr{}, err
		}
		expr = &def.Binary{
			Left:  expr,
			Token: operator,
			Right: right,
		}
	}
	return expr, nil
}

func factor() (def.Expr, error) {
	expr, err := unary()
	if err != nil {
		return &def.EmptyExpr{}, err
	}
	for match(def.SLASH, def.STAR) {
		operator := previous()
		right, errInside := unary()
		if errInside != nil {
			return &def.EmptyExpr{}, err
		}
		expr = &def.Binary{
			Left:  expr,
			Token: operator,
			Right: right,
		}
	}
	return expr, nil
}

func unary() (def.Expr, error) {
	if match(def.BANG, def.MINUS) {
		operator := previous()
		right, err := unary()
		if err != nil {
			return &def.Unary{}, err
		}
		return &def.Unary{
			Token: operator,
			Right: right,
		}, nil
	}
	return primary()
}

func primary() (def.Expr, error) {
	if match(def.FALSE) {
		return &def.Literal{Value: false}, nil
	}

	if match(def.TRUE) {
		return &def.Literal{Value: true}, nil
	}

	if match(def.NIL) {
		return &def.Literal{Value: nil}, nil
	}

	if match(def.NUMBER, def.STRING) {
		return &def.Literal{Value: previous().Literal}, nil
	}

	if match(def.IDENTIFIER) {
		return &def.Variable{Name: previous()}, nil
	}

	if match(def.LEFTPAREN) {
		expr, _ := expression()
		consume(def.RIGHTPAREN, "EXPECT '(' after expression")
		return &def.Grouping{Expression: expr}, nil
	}
	return nil, reportError(peek(), "Expects expression")
}

func synchronize() {
	advance()

	for !isAtEnd() {
		if previous().Type == def.SEMICOLON {
			return
		}
		switch peek().Type {
		case def.CLASS:
		case def.FUN:
		case def.VAR:
		case def.FOR:
		case def.IF:
		case def.WHILE:
		case def.PRINT:
		case def.RETURN:
			return
		}
		advance()
	}
}

func consume(tokenType def.TokenType, message string) (def.Token, error) {
	if check(tokenType) {
		return advance(), nil
	}
	return def.Token{}, reportError(peek(), message)
}

func reportError(token def.Token, message string) error {
	return def.CreateError(token, message)
}

func match(tokenTypes ...def.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if check(tokenType) {
			advance()
			return true
		}
	}
	return false
}

func check(tokenType def.TokenType) bool {
	if isAtEnd() {
		return false
	}
	return tokens[current].Type == tokenType
}

func peek() def.Token {
	return tokens[current]
}

func advance() def.Token {
	if !isAtEnd() {
		current++
	}
	return previous()
}

func isAtEnd() bool {
	return peek().Type == def.EOF
}

func previous() def.Token {
	return tokens[current-1]
}
