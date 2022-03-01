package parser

import (
	"fmt"

	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/token"
)

func (p *parser) parseProgram() *ast.Program {
	program := &ast.Program{}

	for !p.atEnd() {
		program.Statements = append(program.Statements, p.parseStatement())
	}
	return program
}

func (p *parser) parseStatement() ast.Statement {
	var stmt ast.Statement

	switch p.pTok {
	case token.LBRACE:
		stmt = p.parseBlockStmt()
	case token.LET:
		stmt = p.parseLetStmt()
	case token.IF:
		stmt = p.parseIfStatement()
	case token.FOR:
		stmt = p.parseForStmt()
	case token.STRING, token.NOT:
		stmt = p.parseCmdStmt()
	default:
		p.error(p.pPos, fmt.Errorf("illegal token %s at line start", p.pTok))
		p.next()
		return nil
	}

	if !p.match(token.SEMICOLON) {
		p.error(p.pPos, fmt.Errorf("unexpected token %s, expected %s", p.pTok, token.SEMICOLON))
	}

	return stmt
}

func (p *parser) parseBlockStmt() *ast.BlockStatement {
	block := &ast.BlockStatement{}

	if !p.match(token.LBRACE) {
		p.error(p.pPos, fmt.Errorf("expected %s, received %s", token.LBRACE, p.pTok))
	}

	for p.pTok != token.RBRACE && !p.atEnd() {
		block.Statements = append(block.Statements, p.parseStatement())
	}

	if !p.match(token.RBRACE) {
		p.error(p.pPos, fmt.Errorf("expected %s, received %s", token.RBRACE, p.pTok))
	}

	return block
}

func (p *parser) parseLetStmt() *ast.LetStatement {
	if !p.match(token.LET) {
		p.error(p.pPos, fmt.Errorf("expected %s, received %s", token.LET, p.pTok))
	}

	let := &ast.LetStatement{
		Expression: p.parseExprAssign(),
	}

	return let
}

func (p *parser) parseIfStatement() *ast.IfStatement {
	if !p.match(token.IF) {
		p.error(p.pPos, fmt.Errorf("expected %s, received %s", token.IF, p.pTok))
	}

	stmt := ast.IfStatement{
		Condition: p.parseExpression(),
		BlockStmt: p.parseBlockStmt(),
	}

	for p.match(token.ELIF) {
		stmt.ElifBlock = append(stmt.ElifBlock, ast.ElifBlock{
			Condition: p.parseExpression(),
			BlockStmt: p.parseBlockStmt(),
		})
	}

	if p.match(token.ELSE) {
		stmt.ElseBlock = p.parseBlockStmt()
	}

	return &stmt
}

func (p *parser) parseForStmt() *ast.ForStatement {
	if !p.match(token.FOR) {
		p.error(p.pPos, fmt.Errorf("expected %s, received %s", token.FOR, p.pTok))
	}

	var condition ast.Expression

	if !p.check(token.LBRACE) {
		condition = p.parseExpression()
	}

	return &ast.ForStatement{
		Condition: condition,
		BlockStmt: p.parseBlockStmt(),
	}
}

func (p *parser) parseCmdStmt() *ast.CmdStatement {
	return &ast.CmdStatement{
		Command: p.parseCommand(),
	}
}
