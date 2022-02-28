package parser

import (
	"fmt"
	"strconv"

	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/token"
)

func (p *parser) parseExprAssign() ast.Expression {
	expr := p.parseExpression()

	if p.match(token.DEFINE, token.ASSIGN, token.ADD_ASSIGN, token.ADD_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN, token.AND_ASSIGN, token.OR_ASSIGN, token.XOR_ASSIGN, token.SHL_ASSIGN, token.SHR_ASSIGN, token.AND_NOT_ASSIGN) {
		if target, ok := expr.(ast.Assignable); ok {
			return &ast.AssignExpression{
				Left:     target,
				Operator: p.current(),
				Right:    p.parseExpression(),
			}
		}

		p.error(fmt.Errorf("invalid assignment target"))
	}

	return expr
}

func (p *parser) parseExpression() ast.Expression {
	return p.parseExprPrec1()
}

func (p *parser) parseExprPrec1() ast.Expression {
	expr := p.parseExprPrec2()

	for p.match(token.LOR) {
		tok := p.current()
		right := p.parseExprPrec2()
		expr = &ast.LogicalExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr
}

func (p *parser) parseExprPrec2() ast.Expression {
	expr := p.parseExprPrec3()

	for p.match(token.LAND) {
		tok := p.current()
		right := p.parseExprPrec3()
		expr = &ast.LogicalExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr
}

func (p *parser) parseExprPrec3() ast.Expression {
	expr := p.parseExprPrec4()

	for p.match(token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ) {
		tok := p.current()
		right := p.parseExprPrec4()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr
}

func (p *parser) parseExprPrec4() ast.Expression {
	expr := p.parseExprPrec5()

	for p.match(token.ADD, token.SUB, token.OR, token.XOR) {
		tok := p.current()
		right := p.parseExprPrec5()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr
}

func (p *parser) parseExprPrec5() ast.Expression {
	expr := p.parseExprUnary()

	for p.match(token.MUL, token.QUO, token.REM, token.SHL, token.SHR, token.AND, token.AND_NOT) {
		tok := p.current()
		right := p.parseExprUnary()
		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr
}

func (p *parser) parseExprUnary() ast.Expression {
	if p.match(token.ADD, token.SUB, token.XOR, token.NOT) {
		tok := p.current()
		right := p.parseExprLiteral()
		return &ast.UnaryExpression{
			Operator: tok,
			Right:    right,
		}
	}

	return p.parseExprLiteral()
}

func (p *parser) parseExprLiteral() ast.Expression {
	switch {
	case p.match(token.FUNC):
		return &ast.FunctionLiteral{
			Token: p.current(),
			Block: p.parseBlockStmt(),
		}
	case p.match(token.IDENT):
		return &ast.VariableExpression{
			Name: p.current(),
		}
	case p.match(token.FLOAT):
		val, err := strconv.ParseFloat(p.current().Literal, 64)
		if err != nil {
			p.error(err)
		}

		return &ast.NumberLiteral{
			Token: p.current(),
			Value: val,
		}
	case p.match(token.STRING):
		val, err := strconv.Unquote(p.current().Literal)
		if err != nil {
			p.error(err)
		}

		return &ast.StringLiteral{
			Token: p.current(),
			Value: val,
		}
	case p.match(token.LBRACK):
		// TODO: parse arrays and objects
	case p.match(token.LPAREN):
		group := &ast.GroupExpression{
			Right: p.parseExpression(),
		}

		if !p.match(token.RPAREN) {
			p.error(fmt.Errorf("expected %s, received %s", token.RPAREN, p.pTok))
		}

		return group
	}
	return nil
}
