package parser

import (
	"fmt"
	"loxlang/parser/def"
)

var tokens []def.Token
var stmts []def.Stmt
var current int

// Parse is cool
func Parse(input []def.Token) ([]def.Stmt, error) {
	tokens = input
	current = 0
	stmts = []def.Stmt{}
	for !isAtEnd() {
		stmt, _ := declaration()
		if def.HadError {
			return []def.Stmt{}, nil
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func declaration() (def.Stmt, error) {
	if match(def.FUN) {
		funStmt, funErr := function("function")
		if funErr != nil {
			return nil, funErr
		}
		return funStmt, nil
	}
	if match(def.VAR) {
		varStmt, err := varDeclaration()
		if err != nil {
			def.HadError = true
			synchronize()
		}
		return varStmt, nil
	}
	stmt, err := statement()
	if err != nil {
		def.HadError = true
		synchronize()
	}
	return stmt, nil
}

func function(kind string) (def.Stmt, error) {
	name, err := consume(def.IDENTIFIER, fmt.Sprintf("Expected %s name.", kind))
	if err != nil {
		return nil, err
	}
	_, err = consume(def.LEFTPAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
	if err != nil {
		return nil, err
	}
	params := []def.Token{}
	if !check(def.RIGHTPAREN) {
		for {
			if len(params) >= 127 {
				def.CreateError(peek(), "Can't have more than 127 parameters")
			}
			paramID, paramErr := consume(def.IDENTIFIER, "Expect parameter name")
			if paramErr != nil {
				return nil, paramErr
			}
			params = append(params, paramID)
			if !match(def.COMMA) {
				break
			}
		}
	}
	_, err = consume(def.RIGHTPAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}
	_, err = consume(def.LEFTBRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	if err != nil {
		return nil, err
	}
	body, bodyErr := block()
	if bodyErr != nil {
		return nil, bodyErr
	}
	return &def.Function{
		Name:   name,
		Params: params,
		Body:   body,
	}, nil
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
	_, err = consume(def.SEMICOLON, "Expect ';' after variable declaration")
	if err != nil {
		return nil, err
	}
	return &def.Var{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func statement() (def.Stmt, error) {
	if match(def.FOR) {
		return forStatement()
	}

	if match(def.IF) {
		return ifStatement()
	}

	if match(def.PRINT) {
		return printStatement()
	}

	if match(def.RETURN) {
		return returnStatement()
	}

	if match(def.WHILE) {
		return whileStatement()
	}

	if match(def.BREAK) {
		return breakStatement()
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

func breakStatement() (def.Stmt, error) {
	_, err := consume(def.SEMICOLON, "Expect ';' after break keyword.")
	if err != nil {
		return nil, err
	}
	return &def.ControlFlow{
		Type: def.CONTROLFLOWBREAK,
	}, nil
}

func forStatement() (def.Stmt, error) {
	_, err := consume(def.LEFTPAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	var initializer def.Stmt
	if match(def.SEMICOLON) {
		initializer = nil
	} else if match(def.VAR) {
		initializer, err = varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition def.Expr
	if !check(def.SEMICOLON) {
		condition, err = expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = consume(def.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment def.Expr
	if !check(def.RIGHTPAREN) {
		increment, err = expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = consume(def.RIGHTPAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}
	var body def.Stmt
	body, err = statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		stmts := []def.Stmt{}
		body = &def.Block{
			Stmts: append(stmts, body,
				&def.ExprStmt{
					Expr: increment,
				}),
		}
	}

	if condition == nil {
		condition = &def.Literal{Value: true}
	}
	body = &def.While{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		stmts := []def.Stmt{}
		body = &def.Block{
			Stmts: append(stmts, initializer, body),
		}
	}

	return body, nil
}

func returnStatement() (def.Stmt, error) {
	var value def.Expr
	var err error
	keyword := previous()
	if !check(def.SEMICOLON) {
		value, err = expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = consume(def.SEMICOLON, "Expect ';' after return value")
	if err != nil {
		return nil, err
	}
	return &def.Return{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func whileStatement() (def.Stmt, error) {
	_, err := consume(def.LEFTPAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, condErr := expression()
	if condErr != nil {
		return nil, condErr
	}
	_, err = consume(def.RIGHTPAREN, "Expect ')' after 'while' condition.")
	if err != nil {
		return nil, err
	}
	body, bodyErr := statement()
	if bodyErr != nil {
		return nil, bodyErr
	}
	return &def.While{
		Condition: condition,
		Body:      body,
	}, nil
}

func ifStatement() (def.Stmt, error) {
	_, err := consume(def.LEFTPAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, condErr := expression()
	if condErr != nil {
		return nil, condErr
	}
	_, err = consume(def.RIGHTPAREN, "Expect ')' after 'if' condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, thenErr := statement()
	if thenErr != nil {
		return nil, thenErr
	}
	var elseBranch def.Stmt
	if match(def.ELSE) {
		elseBranch, err = statement()
		if err != nil {
			return nil, err
		}
	}
	return &def.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil

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
	expr, err := or()
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

func or() (def.Expr, error) {
	expr, err := and()
	if err != nil {
		return nil, err
	}
	for match(def.OR) {
		operator := previous()
		right, rErr := and()
		if rErr != nil {
			return nil, rErr
		}
		expr = &def.Logical{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}
	return expr, nil
}

func and() (def.Expr, error) {
	expr, err := equality()
	if err != nil {
		return nil, err
	}
	for match(def.AND) {
		operator := previous()
		right, rErr := equality()
		if rErr != nil {
			return nil, rErr
		}
		expr = &def.Logical{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
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
	return call()
}

func call() (def.Expr, error) {
	expr, err := primary()
	if err != nil {
		return nil, err
	}
	for {
		if match(def.LEFTPAREN) {
			expr, err = finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func finishCall(callee def.Expr) (def.Expr, error) {
	args := []def.Expr{}
	if !check(def.RIGHTPAREN) {
		for {
			if len(args) >= 127 {
				return nil, reportError(peek(), "Can't have more than 127 arguments in a function")
			}
			expr, err := expression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
			if !match(def.COMMA) {
				break
			}
		}
	}
	paren, parenErr := consume(def.RIGHTPAREN, "Expect ')' after argument list")
	if parenErr != nil {
		return nil, parenErr
	}
	return &def.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: args,
	}, nil

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
