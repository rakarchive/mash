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

	p.consume(token.SEMICOLON)
	return stmt
}

func (p *parser) parseBlockStmt() *ast.BlockStatement {
	block := &ast.BlockStatement{}

	p.consume(token.LBRACE)

	for p.pTok != token.RBRACE && !p.atEnd() {
		block.Statements = append(block.Statements, p.parseStatement())
	}

	p.consume(token.RBRACE)
	return block
}

func (p *parser) parseLetStmt() *ast.LetStatement {
	p.consume(token.LET)

	let := &ast.LetStatement{
		Expression: p.parseExprAssign(),
	}

	return let
}

func (p *parser) parseIfStatement() *ast.IfStatement {
	p.consume(token.IF)

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
	p.consume(token.FOR)

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
