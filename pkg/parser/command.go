package parser

import (
	"fmt"

	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/token"
)

func (p *parser) parseCommand() (ast.Command, error) {
	return p.parseCmdLor()
}

func (p *parser) parseCmdLor() (ast.Command, error) {
	expr, err := p.parseCmdAnd()
	if err != nil {
		return nil, err
	}

	for p.match(token.LOR) {
		tok := p.current()
		right, err := p.parseCmdAnd()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalCommand{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) parseCmdAnd() (ast.Command, error) {
	expr, err := p.parseCmdNot()
	if err != nil {
		return nil, err
	}

	for p.match(token.LAND) {
		tok := p.current()
		right, err := p.parseCmdNot()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalCommand{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) parseCmdNot() (ast.Command, error) {
	if p.match(token.NOT) {
		tok := p.current()
		right, err := p.parseCmdPipe()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryCommand{
			Operator: tok,
			Right:    right,
		}, nil
	}

	return p.parseCmdPipe()
}

func (p *parser) parseCmdPipe() (ast.Command, error) {
	expr, err := p.parseCmdLit()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		tok := p.current()
		right, err := p.parseCmdLit()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryCommand{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) parseCmdLit() (ast.Command, error) {
	if !p.match(token.STRING) {
		return nil, fmt.Errorf("unexpected token %s", p.pTok)
	}

	lit := &ast.LiteralCommand{
		Cmd: p.current(),
	}

	for p.match(token.STRING) {
		lit.Args = append(lit.Args, p.current())
	}

	return lit, nil
}
