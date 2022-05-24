// Copyright © 2022 Rak Laptudirm <raklaptudirm@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"fmt"
	"strconv"

	"laptudirm.com/x/mash/pkg/ast"
	"laptudirm.com/x/mash/pkg/token"
)

// AssignExpression = Assignable assign_op Expression .
func (p *parser) parseAssignExpression() (ast.Expression, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.match(token.Define, token.Assign, token.AdditionAssign, token.AdditionAssign, token.MultiplicationAssign, token.QuotientAssign, token.RemainderAssign, token.AndAssign, token.OrAssign, token.XorAssign, token.ShiftLeftAssign, token.ShiftRightAssign, token.AndNotAssign) {
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

// Expression = OrExpression .
func (p *parser) parseExpression() (ast.Expression, error) {
	return p.parseOrExpression()
}

// OrExpression = AndExpression { "||" OrExpression } .
func (p *parser) parseOrExpression() (ast.Expression, error) {
	expr, err := p.parseAndExpression()
	if err != nil {
		return nil, err
	}

	for p.match(token.LogicalOr) {
		tok := p.current()
		right, err := p.parseAndExpression()
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

// AndExpression = RelExpression { "&&" AndExpression } .
func (p *parser) parseAndExpression() (ast.Expression, error) {
	expr, err := p.parseRelExpression()
	if err != nil {
		return nil, err
	}

	for p.match(token.LogicalAnd) {
		tok := p.current()
		right, err := p.parseRelExpression()
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

// RelExpression = AddExpression { rel_op RelExpression } .
func (p *parser) parseRelExpression() (ast.Expression, error) {
	expr, err := p.parseAddExpression()
	if err != nil {
		return nil, err
	}

	for p.match(token.Equal, token.NotEqual, token.LessThan, token.LessThanEqual, token.GreaterThan, token.GreaterThanEqual) {
		tok := p.current()
		right, err := p.parseAddExpression()
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

// AddExpression = MulExpression { add_op AddExpression } .
func (p *parser) parseAddExpression() (ast.Expression, error) {
	expr, err := p.parseMulExpression()
	if err != nil {
		return nil, err
	}

	for p.match(token.Addition, token.Subtraction, token.Or, token.Xor) {
		tok := p.current()
		right, err := p.parseMulExpression()
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

// MulExpression = UnaryExpression { mul_op MulExpression } .
func (p *parser) parseMulExpression() (ast.Expression, error) {
	expr, err := p.parseUnaryExpression()
	if err != nil {
		return nil, err
	}

	for p.match(token.Multiplication, token.Quotient, token.Remainder, token.ShiftLeft, token.ShiftRight, token.And, token.AndNot) {
		tok := p.current()
		right, err := p.parseUnaryExpression()
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

// UnaryExpression = PrimaryExpression | unary_op UnaryExpression .
func (p *parser) parseUnaryExpression() (ast.Expression, error) {
	if p.match(token.Addition, token.Subtraction, token.Xor, token.Not) {
		tok := p.current()
		right, err := p.parsePrimaryExpression()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpression{
			Operator: tok,
			Right:    right,
		}, nil
	}

	return p.parsePrimaryExpression()
}

// PrimaryExpression = Operand { Selector | Index | Arguments } .
func (p *parser) parsePrimaryExpression() (ast.Expression, error) {
	expr, err := p.parseOperand()
	if err != nil {
		return nil, err
	}

	for {
		var parseFunc func(ast.Expression) (ast.Expression, error)

		switch p.pTok {
		case token.Period:
			parseFunc = p.parseSelector
		case token.LeftBrack:
			parseFunc = p.parseIndex
		case token.LeftParen:
			parseFunc = p.parseArguments
		default:
			return expr, nil
		}

		new, err := parseFunc(expr)
		if err != nil {
			return nil, err
		}

		expr = new
	}
}

// Selector  = "." identifier .
func (p *parser) parseSelector(expr ast.Expression) (ast.Expression, error) {
	p.match(token.Period)
	if !p.match(token.Identifier) {
		return nil, fmt.Errorf("expected identifier, received %s", p.pTok)
	}

	return &ast.SelectorExpression{
		Name:  expr,
		Index: p.current(),
	}, nil
}

// Index = "[" Expression "]" .
func (p *parser) parseIndex(expr ast.Expression) (ast.Expression, error) {
	p.match(token.LeftBrack)

	name, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.match(token.RightBrack) {
		return nil, fmt.Errorf("expected '}', received %s", p.pTok)
	}

	return &ast.GetExpression{
		Name: name,
		Expr: expr,
	}, nil
}

// Arguments = "(" ExpressionList ")" .
func (p *parser) parseArguments(expr ast.Expression) (ast.Expression, error) {
	var args []ast.Expression

	p.match(token.LeftParen)
	paren := p.current()

	args, err := p.parseExpressionList(token.RightParen)
	if err != nil {
		return nil, err
	}

	return &ast.CallExpression{
		Callee:      expr,
		Parenthesis: paren,
		Arguments:   args,
	}, nil
}

// Operand = Literal | "(" Expression ")" .
func (p *parser) parseOperand() (ast.Expression, error) {
	switch {
	case p.match(token.LeftParen):
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if !p.match(token.RightParen) {
			return nil, fmt.Errorf("expected ')', received %s", p.pTok)
		}

		return expr, nil
	default:
		return p.parseLiteral()
	}
}

// Literal = BasicLit | CompositeLit | FunctionLit .
func (p *parser) parseLiteral() (ast.Expression, error) {
	switch p.pTok {
	case token.Identifier, token.Number, token.String:
		return p.parseBasicLit()
	case token.LeftBrack:
		return p.parseArrayLit()
	case token.Obj:
		return p.parseObjectLit()
	case token.Func:
		return p.parseFunctionLit()
	case token.Template:
		return p.parseTemplateLit()
	default:
		return nil, fmt.Errorf("invalid literal %s", p.pTok)
	}
}

// BasicLit = identifier | number_lit | string_lit .
func (p *parser) parseBasicLit() (ast.Expression, error) {
	switch p.next(); p.tok {
	case token.Identifier:
		return &ast.VariableExpression{
			Name: p.current(),
		}, nil
	case token.Number:
		val, err := strconv.ParseFloat(p.lit, 64)
		if err != nil {
			return nil, err
		}

		return &ast.NumberLiteral{
			Token: p.current(),
			Value: val,
		}, nil
	case token.String:
		val, err := strconv.Unquote(p.lit)
		if err != nil {
			return nil, err
		}

		return &ast.StringLiteral{
			Token: p.current(),
			Value: val,
		}, nil
	default:
		panic("invalid token provided to BasicLit")
	}
}

// ArrayLit = "[" ExpressionList "]" .
func (p *parser) parseArrayLit() (*ast.ArrayLiteral, error) {
	p.match(token.LeftBrack)
	brack := p.current()

	list, err := p.parseExpressionList(token.RightBrack)
	if err != nil {
		return nil, err
	}

	return &ast.ArrayLiteral{
		Token:    brack,
		Elements: list,
	}, nil
}

// ObjectLit = "obj" "[" ObjectEntryList [ "," ] "]" .
func (p *parser) parseObjectLit() (*ast.ObjectLiteral, error) {
	p.match(token.Obj)
	obj := p.current()

	if !p.match(token.LeftBrack) {
		return nil, fmt.Errorf("expected '[', received %s", p.pTok)
	}

	elements := make(map[ast.Expression]ast.Expression)
	for !p.check(token.RightBrack) && !p.atEnd() {
		key, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if !p.match(token.Colon) {
			return nil, fmt.Errorf("expected ':', received %s", p.pTok)
		}

		value, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		elements[key] = value

		if !p.match(token.Comma) && !p.check(token.RightBrack) {
			return nil, fmt.Errorf("expected ']', received %s", p.pTok)
		}
	}

	if !p.match(token.RightBrack) {
		return nil, fmt.Errorf("expected ']', received %s", p.pTok)
	}

	return &ast.ObjectLiteral{
		Token:    obj,
		Elements: elements,
	}, nil
}

// FunctionLit = "func" Block .
func (p *parser) parseFunctionLit() (*ast.FunctionLiteral, error) {
	p.match(token.Func)
	tok := p.current()

	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionLiteral{
		Token: tok,
		Block: block,
	}, nil
}

// TemplateLit = "'" _embedded_string_val "'" .
func (p *parser) parseTemplateLit() (*ast.TemplateLiteral, error) {
	p.match(token.Template)

	var components []token.Token
	var expressions []ast.Expression

	for {
		p.match(token.String)
		components = append(components, p.current())

		if p.match(token.Template) {
			break
		}

		p.match(token.LeftBrace)

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)

		p.match(token.RightBrace)
	}

	return &ast.TemplateLiteral{
		Expressions: expressions,
		Components:  components,
	}, nil
}

// ExpressionList = Expression { "," Expression } .
func (p *parser) parseExpressionList(eol token.Type) ([]ast.Expression, error) {
	var list []ast.Expression

	for !p.check(eol) && !p.atEnd() {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		list = append(list, expr)

		if !p.match(token.Comma) && !p.check(eol) {
			return nil, fmt.Errorf("expected %s, received %s", eol, p.pTok)
		}
	}

	if !p.match(eol) {
		return nil, fmt.Errorf("expected %s, received EOF", eol)
	}

	return list, nil
}
