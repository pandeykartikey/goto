package parser

import (
	"fmt"
	"strconv"

	"github.com/pandeykartikey/goto/ast"
	"github.com/pandeykartikey/goto/lexer"
	"github.com/pandeykartikey/goto/token"
)

const ( // These represent the operator precedence values.
	_int = iota
	LOWEST
	LOGICAL     // && or ||
	EQUALS      // ==
	LESSGREATER // > or <
	PLUS        // +
	MULTIPLY    // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // []
)

var precedences = map[token.Type]int{
	token.AND:      LOGICAL,
	token.OR:       LOGICAL,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LT_EQ:    LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.PLUS:     PLUS,
	token.MINUS:    PLUS,
	token.DIVIDE:   MULTIPLY,
	token.MULTIPLY: MULTIPLY,
	token.MOD:      MULTIPLY,
	token.POW:      MULTIPLY,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParsefn func() ast.Expression
	infixParsefn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string

	prefixParsefns map[token.Type]prefixParsefn
	infixParsefns  map[token.Type]infixParsefn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParsefns = make(map[token.Type]prefixParsefn)
	prefixfns := []struct {
		token   token.Type
		parsefn prefixParsefn
	}{
		{token.IDENT, p.parseIdentifier},
		{token.INT, p.parseIntegerLiteral},
		{token.NOT, p.parsePrefixExpression},
		{token.MINUS, p.parsePrefixExpression},
		{token.TRUE, p.parseBoolean},
		{token.FALSE, p.parseBoolean},
		{token.STRING, p.parseString},
		{token.LPAREN, p.parseGroupedExpression},
		{token.LBRACKET, p.parseList},
	}

	for _, fn := range prefixfns {
		p.registerPrefix(fn.token, fn.parsefn)
	}

	p.infixParsefns = make(map[token.Type]infixParsefn)
	for keys := range precedences {
		p.registerInfix(keys, p.parseInfixExpression)
	}

	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.setToken() // Only to be called for initialization of Parser pointers

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken(count ...int64) {
	n := int64(1)
	if count != nil {
		if len(count) != 1 {
			panic("nextToken takes only one or no parameters")
		}
		n = count[0]
	}
	for i := int64(0); i < n; i++ {
		p.currToken = p.peekToken
		p.peekToken = p.l.NextToken()
	}
}

func (p *Parser) setToken() {
	p.currToken = p.l.NextToken()
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currTokenIs(t token.Type) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) tokenError(exp token.Type, t token.Type) {
	msg := fmt.Sprintf("expected token to be %s , got %s instead", exp, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectCurr(t token.Type) bool {

	if p.currTokenIs(t) {
		return true
	}

	p.tokenError(t, p.currToken.Type)
	return false
}

func (p *Parser) expectPeek(t token.Type) bool {

	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.tokenError(t, p.peekToken.Type)
	return false
}

func (p *Parser) registerPrefix(Type token.Type, fn prefixParsefn) {
	p.prefixParsefns[Type] = fn
}

func (p *Parser) registerInfix(Type token.Type, fn infixParsefn) {
	p.infixParsefns[Type] = fn
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseIdentifier() ast.Expression {
	if !p.expectCurr(token.IDENT) {
		return nil
	}
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixexp := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}

	p.nextToken()

	prefixexp.Right = p.parseExpression(PREFIX)

	return prefixexp
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infixexp := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	precedence := p.currPrecedence()
	p.nextToken()

	infixexp.Right = p.parseExpression(precedence)
	return infixexp
}

// parses exprA, ... ,exprZ Initial currtoken at exprA and Final after exprZ
func (p *Parser) parseExpressionList() *ast.ExpressionList {
	args := &ast.ExpressionList{Token: p.currToken}

	for !p.currTokenIs(token.EOF) && !p.currTokenIs(token.RPAREN) {
		exp := p.parseExpression(LOWEST)
		args.Expressions = append(args.Expressions, &exp)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken(2)
			continue
		}

		p.nextToken()
		return args
	}

	if p.currTokenIs(token.EOF) {
		p.errors = append(p.errors, "End Of File encountered while parsing")
	}

	return nil
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currToken}
	fname, ok := left.(*ast.Identifier)
	if !ok {
		return nil
	}
	exp.FunctionName = fname

	p.nextToken()
	exp.ArgumentList = p.parseExpressionList()

	if !p.expectCurr(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseExpression(precedence int) ast.Expression { // returns expression on the same or higher precedence level
	prefix := p.prefixParsefns[p.currToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMI) && precedence < p.peekPrecedence() {
		infix := p.infixParsefns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if p.expectPeek(token.RPAREN) {
		return exp
	}

	return nil
}

// isExpression is used to maintain difference between Assignment statement and expression
func (p *Parser) parseAssignment(isExpression bool) *ast.Assignment {
	assign := &ast.Assignment{IsExpression: isExpression}

	if p.currTokenIs(token.VAR) {
		assign.Token = p.currToken
		p.nextToken()
	}

	assign.NameList = p.parseIdentifierList()

	if assign.Token.Type == token.VAR && !isExpression && p.currTokenIs(token.SEMI) {
		return assign
	}

	if !p.expectCurr(token.ASSIGN) {
		return nil
	}

	if assign.Token.Type != token.VAR {
		assign.Token = p.currToken
	}

	p.nextToken()

	assign.ValueList = p.parseExpressionList()

	if assign.ValueList == nil || assign.NameList == nil || len(assign.ValueList.Expressions) != len(assign.NameList.Identifiers) {
		p.errors = append(p.errors, "Mismatch in number of values on both side of =")
	}

	if !isExpression && !p.expectCurr(token.SEMI) {
		return nil
	}

	return assign
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMI) {
		return nil
	}

	return stmt
}

func (p *Parser) parseLoopControlStatement() *ast.LoopControlStatement {
	stmt := &ast.LoopControlStatement{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.SEMI) {
		return nil
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	p.nextToken()
	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()

	}

	if !p.expectCurr(token.RBRACE) {
		return nil
	}

	return block
}

func (p *Parser) parseList() ast.Expression {
	list := &ast.List{Token: p.currToken}
	p.nextToken()
	list.Elements = p.parseExpressionList()

	if !p.expectCurr(token.RBRACKET) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.currToken, Left: left}
	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currToken}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if !p.peekTokenIs(token.ELSE) {
		return stmt
	}

	p.nextToken()

	if p.peekTokenIs(token.IF) {
		p.nextToken()
		stmt.FollowIf = p.parseIfStatement()
	} else if !p.expectPeek(token.LBRACE) {
		return nil
	} else {
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.currToken}

	p.nextToken()

	if !p.currTokenIs(token.SEMI) {
		stmt.Init = p.parseAssignment(true)
	}

	if !p.expectCurr(token.SEMI) {
		return nil
	}
	p.nextToken()

	if !p.currTokenIs(token.SEMI) {
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.SEMI) {
		return nil
	}

	p.nextToken()
	if !p.currTokenIs(token.LBRACE) {
		stmt.Update = p.parseAssignment(true)
	}

	if !p.expectCurr(token.LBRACE) {
		return nil
	}

	stmt.ForBody = p.parseBlockStatement()

	return stmt
}

// parses identA, ... ,identZ Initial currtoken at identA and Final after identZ
func (p *Parser) parseIdentifierList() *ast.IdentifierList {
	identlist := &ast.IdentifierList{Token: p.currToken}

	for !p.currTokenIs(token.EOF) && p.currTokenIs(token.IDENT) {

		ident, ok := p.parseIdentifier().(*ast.Identifier)

		if !ok {
			return nil
		}

		if ident != nil {
			identlist.Identifiers = append(identlist.Identifiers, ident)
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken(2)
			continue
		}

		p.nextToken()
		return identlist
	}

	if p.currTokenIs(token.EOF) {
		p.errors = append(p.errors, "End Of File encountered while parsing")
	}

	return nil
}

func (p *Parser) parseFuncStatement() *ast.FuncStatement {
	stmt := &ast.FuncStatement{Token: p.currToken}

	p.nextToken()

	name, ok := p.parseIdentifier().(*ast.Identifier)

	if !ok {
		return nil
	}

	stmt.Name = name

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	stmt.ParameterList = p.parseIdentifierList()

	if !p.expectCurr(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.FuncBody = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMI) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.VAR:
		return p.parseAssignment(false)
	case token.RETURN:
		return p.parseReturnStatement()
	case token.BREAK, token.CONTINUE:
		return p.parseLoopControlStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	case token.FUNC:
		return p.parseFuncStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.ILLEGAL:
		p.errors = append(p.errors, "ILLEGAL Token encountered")
		return nil
	case token.SEMI:
		return nil
	case token.EOF:
		return nil
	case token.IDENT:
		if p.peekTokenIs(token.COMMA) || p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignment(false)
		}
		fallthrough
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) PrintParseErrors() {
	for _, msg := range p.errors {
		fmt.Println("Error: ", msg)
	}
}
