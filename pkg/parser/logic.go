package parser

import (
	"fmt"

	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/token"
)

func (p *parser) parseProgram() *ast.Program {
	program := &ast.Program{}

	for p.next(); p.tok != token.EOF; p.next() {
		program.Statements = append(program.Statements, p.parseStatement())
	}
	return program
}

func (p *parser) parseStatement() ast.Statement {
	switch p.tok {
	case token.LBRACE:
		return p.parseBlock()
	case token.LET:
		// parse expression
	case token.IF:
		// parse if
	case token.FOR:
		// parse for
	case token.STRING:
		// parse command
	default:
		p.error(fmt.Errorf("illegal token %#v at line start", p.tok.String()))
	}

	return nil
}

func (p *parser) parseBlock() *ast.BlockStatement {
	block := &ast.BlockStatement{}

	for p.next(); p.tok != token.RBRACE; p.next() {
		block.Statements = append(block.Statements, p.parseStatement())
	}
	return block
}