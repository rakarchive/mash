package parser

import (
	"fmt"
	"strconv"

	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/token"
)

func (p *parser) parseExprAssign() (ast.Expression, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.match(token.DEFINE, token.ASSIGN, token.ADD_ASSIGN, token.ADD_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN, token.REM_ASSIGN, token.AND_ASSIGN, token.OR_ASSIGN, token.XOR_ASSIGN, token.SHL_ASSIGN, token.SHR_ASSIGN, token.AND_NOT_ASSIGN) {
		if target, ok := expr.(ast.Assignable); ok {
			tok := p.current()
			right, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			return &ast.AssignExpression{
				Left:     target,
				Operator: tok,
				Right:    right,
			}, nil
		}

		return nil, fmt.Errorf("invalid assignment target")
	}

	return expr, nil
}

func (p *parser) parseExpression() (ast.Expression, error) {
	return p.parseExprPrec1()
}

func (p *parser) parseExprPrec1() (ast.Expression, error) {
	expr, err := p.parseExprPrec2()
	if err != nil {
		return nil, err
	}

	for p.match(token.LOR) {
		tok := p.current()
		right, err := p.parseExprPrec2()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) parseExprPrec2() (ast.Expression, error) {
	expr, err := p.parseExprPrec3()
	if err != nil {
		return nil, err
	}

	for p.match(token.LAND) {
		tok := p.current()
		right, err := p.parseExprPrec3()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) parseExprPrec3() (ast.Expression, error) {
	expr, err := p.parseExprPrec4()
	if err != nil {
		return nil, err
	}

	for p.match(token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ) {
		tok := p.current()
		right, err := p.parseExprPrec4()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) parseExprPrec4() (ast.Expression, error) {
	expr, err := p.parseExprPrec5()
	if err != nil {
		return nil, err
	}

	for p.match(token.ADD, token.SUB, token.OR, token.XOR) {
		tok := p.current()
		right, err := p.parseExprPrec5()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) parseExprPrec5() (ast.Expression, error) {
	expr, err := p.parseExprUnary()
	if err != nil {
		return nil, err
	}

	for p.match(token.MUL, token.QUO, token.REM, token.SHL, token.SHR, token.AND, token.AND_NOT) {
		tok := p.current()
		right, err := p.parseExprUnary()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpression{
			Left:     expr,
			Operator: tok,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) parseExprUnary() (ast.Expression, error) {
	if p.match(token.ADD, token.SUB, token.XOR, token.NOT) {
		tok := p.current()
		right, err := p.parseExprCall()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpression{
			Operator: tok,
			Right:    right,
		}, nil
	}

	return p.parseExprCall()
}

func (p *parser) parseExprCall() (ast.Expression, error) {
	expr, err := p.parseExprLiteral()
	if err != nil {
		return nil, err
	}

	for {
		switch {
		case p.match(token.LBRACK):
			name, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			expr = &ast.GetExpression{
				Name: name,
				Expr: expr,
			}

			if !p.match(token.RBRACK) {
				return nil, fmt.Errorf("expected '}', received %s", p.pTok)
			}
		case p.match(token.LPAREN):
			args := []ast.Expression{}
			tok := p.current()

			for !p.match(token.RPAREN) {
				exp, err := p.parseExpression()
				if err != nil {
					return nil, err
				}

				args = append(args, exp)

				if !p.match(token.COMMA) {
					if !p.match(token.RPAREN) {
						return nil, fmt.Errorf("expected ')', received %s", p.pTok)
					}

					break
				}
			}

			expr = &ast.CallExpression{
				Callee:      expr,
				Parenthesis: tok,
				Arguments:   args,
			}
		default:
			return expr, nil
		}
	}
}

func (p *parser) parseExprLiteral() (ast.Expression, error) {
	switch {
	case p.match(token.FUNC):
		block, err := p.parseBlockStmt()
		if err != nil {
			return nil, err
		}

		return &ast.FunctionLiteral{
			Token: p.current(),
			Block: block,
		}, nil
	case p.match(token.IDENT):
		return &ast.VariableExpression{
			Name: p.current(),
		}, nil
	case p.match(token.FLOAT):
		val, err := strconv.ParseFloat(p.lit, 64)
		if err != nil {
			return nil, err
		}

		return &ast.NumberLiteral{
			Token: p.current(),
			Value: val,
		}, nil
	case p.match(token.STRING):
		val, err := strconv.Unquote(p.lit)
		if err != nil {
			return nil, err
		}

		return &ast.StringLiteral{
			Token: p.current(),
			Value: val,
		}, nil
	case p.match(token.OBJ):
		lit := &ast.ObjectLiteral{
			Token:    p.current(),
			Elements: make(map[ast.Expression]ast.Expression),
		}

		if !p.match(token.LBRACK) {
			return nil, fmt.Errorf("expected '[', received %s", p.pTok)
		}

		if p.match(token.RBRACK) {
			return lit, nil
		}

		for {
			index, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			if !p.match(token.COLON) {
				return nil, fmt.Errorf("expected ':', received %s", p.pTok)
			}

			lit.Elements[index], err = p.parseExpression()
			if err != nil {
				return nil, err
			}

			if p.match(token.COMMA) {
				if p.match(token.RBRACK) {
					break
				}

				continue
			}

			if !p.match(token.RBRACK) {
				return nil, fmt.Errorf("expected ']', received %s", p.pTok)
			}

			break
		}

		return lit, nil
	case p.match(token.LBRACK):
		lit := &ast.ArrayLiteral{
			Token:    p.current(),
			Elements: []ast.Expression{},
		}

		if p.match(token.RBRACK) {
			return lit, nil
		}

		for {
			exp, err := p.parseExpression()
			if err != nil {
				return nil, err
			}

			lit.Elements = append(lit.Elements, exp)

			if p.match(token.COMMA) {
				if p.match(token.RBRACK) {
					break
				}

				continue
			}

			if !p.match(token.RBRACK) {
				return nil, fmt.Errorf("expected ']', received %s", p.pTok)
			}
			break
		}

		return lit, nil
	case p.match(token.LPAREN):
		right, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		group := &ast.GroupExpression{
			Right: right,
		}

		if !p.match(token.RPAREN) {
			return nil, fmt.Errorf("expected ')', received %s", p.pTok)
		}

		return group, nil
	default:
		return nil, fmt.Errorf("invalid literal start %s", p.pTok)
	}
}
