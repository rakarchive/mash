// Copyright © 2021 Rak Laptudirm <raklaptudirm@gmail.com>
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

	"laptudirm.com/x/mash/pkg/ast"
	"laptudirm.com/x/mash/pkg/token"
)

// OrCommand = AndCommand { "||" OrCommand } .
func (p *parser) parseOrCommand() (ast.Command, error) {
	expr, err := p.parseAndCommand()
	if err != nil {
		return nil, err
	}

	for p.match(token.LogicalOr) {
		tok := p.current()
		right, err := p.parseAndCommand()
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

// AndCommand = NotCommand { "&&" AndCommand } .
func (p *parser) parseAndCommand() (ast.Command, error) {
	expr, err := p.parseNotCommand()
	if err != nil {
		return nil, err
	}

	for p.match(token.LogicalAnd) {
		tok := p.current()
		right, err := p.parseNotCommand()
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

// NotCommand = [ "!" ] PipeCommand .
func (p *parser) parseNotCommand() (ast.Command, error) {
	if p.match(token.Not) {
		tok := p.current()
		right, err := p.parsePipeCommand()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryCommand{
			Operator: tok,
			Right:    right,
		}, nil
	}

	return p.parsePipeCommand()
}

// PipeCommand = PrimaryCommand { "|" PipeCommand } .
func (p *parser) parsePipeCommand() (ast.Command, error) {
	expr, err := p.parsePrimaryCommand()
	if err != nil {
		return nil, err
	}

	for p.match(token.Or) {
		tok := p.current()
		right, err := p.parsePrimaryCommand()
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

// PrimaryCommand = command_arg { command_arg } .
func (p *parser) parsePrimaryCommand() (ast.Command, error) {
	if !p.check(token.String, token.Template) {
		return nil, fmt.Errorf("unexpected token %s", p.pTok)
	}

	var components []ast.CommandComponent

componentLoop:
	for {
		var component ast.CommandComponent

		switch p.pTok {
		case token.String:
			p.next()

			component = &ast.StringLiteral{
				Token: p.current(),
				Value: p.lit,
			}
		case token.Template:
			template, err := p.parseTemplateLit()
			if err != nil {
				return nil, err
			}

			component = template
		default:
			break componentLoop
		}

		components = append(components, component)
	}

	return &ast.LiteralCommand{
		Components: components,
	}, nil
}
