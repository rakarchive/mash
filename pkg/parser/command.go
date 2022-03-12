package parser

import (
	"fmt"

	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/token"
)

func (p *parser) parseCommand() ast.Command {
	return p.parseCmdLor()
}

func (p *parser) parseCmdLor() ast.Command {
	expr := p.parseCmdAnd()

	for p.match(token.LOR) {
		tok := p.current()
		right := p.parseCmdAnd()
		expr = &ast.LogicalCommand{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr
}

func (p *parser) parseCmdAnd() ast.Command {
	expr := p.parseCmdNot()

	for p.match(token.LAND) {
		tok := p.current()
		right := p.parseCmdNot()
		expr = &ast.LogicalCommand{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr
}

func (p *parser) parseCmdNot() ast.Command {
	if p.match(token.NOT) {
		tok := p.current()
		right := p.parseCmdPipe()
		return &ast.UnaryCommand{
			Operator: tok,
			Right:    right,
		}
	}

	return p.parseCmdPipe()
}

func (p *parser) parseCmdPipe() ast.Command {
	expr := p.parseCmdLit()

	for p.match(token.OR) {
		tok := p.current()
		right := p.parseCmdLit()
		expr = &ast.BinaryCommand{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr
}

func (p *parser) parseCmdLit() ast.Command {
	if !p.match(token.STRING) {
		p.error(p.pPos, fmt.Errorf("unexpected toke %s", p.pTok))
	}

	lit := &ast.LiteralCommand{
		Cmd: p.current(),
	}

	for p.match(token.STRING) {
		lit.Args = append(lit.Args, p.current())
	}

	return lit
}
