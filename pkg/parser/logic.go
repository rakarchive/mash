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
		// parse expression
	case token.IF:
		// parse if
	case token.FOR:
		// parse for
	case token.STRING, token.NOT:
		stmt = p.parseCmdStmt()
	default:
		p.error(fmt.Errorf("illegal token %s at line start", p.pTok))
		p.next()
		return nil
	}

	if !p.match(token.SEMICOLON) {
		p.error(fmt.Errorf("unexpected token %s, expected %s", p.pTok, token.SEMICOLON))
	}

	return stmt
}

func (p *parser) parseBlockStmt() *ast.BlockStatement {
	block := &ast.BlockStatement{}

	for p.next(); p.tok != token.RBRACE; p.next() {
		block.Statements = append(block.Statements, p.parseStatement())
	}
	return block
}

func (p *parser) parseCmdStmt() *ast.CmdStatement {
	return &ast.CmdStatement{
		Command: p.parseCommand(),
	}
}
