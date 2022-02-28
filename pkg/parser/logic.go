package parser

import (
	"fmt"

	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/token"
)

func (p *parser) parseProgram() *ast.Program {
	program := &ast.Program{}

	for p.next(); p.curToken.Type != token.EOF; p.next() {
		program.Statements = append(program.Statements, p.parseStatement())
	}
	return program
}

func (p *parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LBRACE:
		// parse block
	case token.LET:
		// parse expression
	case token.IF:
		// parse if
	case token.FOR:
		// parse for
	case token.STRING:
		// parse command
	default:
		p.error(fmt.Errorf("illegal token %#v at line start", p.curToken.Type.String()))
	}

	return nil
}
