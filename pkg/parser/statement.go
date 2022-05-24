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

	"laptudirm.com/x/mash/pkg/ast"
	"laptudirm.com/x/mash/pkg/token"
)

// Program = StatementList .
func (p *parser) parseProgram() *ast.Program {
	return &ast.Program{
		Statements: p.parseStatementList(token.Eof),
	}
}

// Block = "{" StatementList "}" .
func (p *parser) parseBlock() (*ast.BlockStatement, error) {
	if !p.match(token.LeftBrace) {
		return nil, fmt.Errorf("expected '{', received %s", p.pTok)
	}

	statements := p.parseStatementList(token.RightBrace)
	if !p.match(token.RightBrace) {
		return nil, fmt.Errorf("expected '}', received %s", p.pTok)
	}

	return &ast.BlockStatement{
		Statements: statements,
	}, nil
}

// StatementList = { Statement } .
func (p *parser) parseStatementList(eos token.Type) []ast.Statement {
	var statements []ast.Statement

	for p.pTok != eos && !p.atEnd() {
		stmt, err := p.parseStatement()
		if err != nil {
			p.error(p.pPos, err)
			// sync parser to avoid cascading errors
			p.synchronize()
			continue
		}

		statements = append(statements, stmt)
	}

	return statements
}

// Statement = ( LetStatement | ForStatement | IfStatement | Block | CommandStatement ) ";" .
func (p *parser) parseStatement() (ast.Statement, error) {
	var stmt ast.Statement
	var err error

	switch p.pTok {
	case token.Let:
		stmt, err = p.parseLetStatement()
	case token.For:
		stmt, err = p.parseForStatement()
	case token.If:
		stmt, err = p.parseIfStatement()
	case token.LeftBrace:
		stmt, err = p.parseBlock()
	case token.String, token.Not:
		stmt, err = p.parseCommandStatement()
	default:
		return nil, fmt.Errorf("illegal token %s at line start", p.pTok)
	}

	// only check for semicolons if no errors have occurred
	if err == nil && !p.match(token.Semicolon) {
		return nil, fmt.Errorf("expected ';', received %s", p.pTok)
	}
	return stmt, err
}

// LetStatement = "let" AssignExpression .
func (p *parser) parseLetStatement() (*ast.LetStatement, error) {
	if !p.match(token.Let) {
		return nil, fmt.Errorf("expected 'let', received %s", p.pTok)
	}

	expr, err := p.parseAssignExpression()
	if err != nil {
		return nil, err
	}

	return &ast.LetStatement{
		Expression: expr,
	}, nil
}

// ForStatement = "for" [ Expression ] Block .
func (p *parser) parseForStatement() (*ast.ForStatement, error) {
	if !p.match(token.For) {
		return nil, fmt.Errorf("expected 'for', received %s", p.pTok)
	}

	var condition ast.Expression
	var err error

	if !p.check(token.LeftBrace) {
		condition, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.ForStatement{
		Condition: condition,
		BlockStmt: block,
	}, nil
}

// IfStatement = "if" Expression Block [ "else" ( IfStatement | Block ) ] .
func (p *parser) parseIfStatement() (*ast.IfStatement, error) {
	if !p.match(token.If) {
		return nil, fmt.Errorf("expected 'if', received %s", p.pTok)
	}

	cond, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	var elseBlock ast.Statement
	if p.match(token.Else) {
		switch p.pTok {
		case token.If:
			stmt, err := p.parseIfStatement()
			if err != nil {
				return nil, err
			}

			elseBlock = stmt
		case token.LeftBrace:
			block, err := p.parseBlock()
			if err != nil {
				return nil, err
			}

			elseBlock = block
		default:
			return nil, fmt.Errorf("illegal token %s after 'else'", p.pTok)
		}
	}

	return &ast.IfStatement{
		Condition: cond,
		BlockStmt: block,
		ElseBlock: elseBlock,
	}, nil
}

// CommandStatement = OrCommand .
func (p *parser) parseCommandStatement() (*ast.CmdStatement, error) {
	cmd, err := p.parseOrCommand()
	if err != nil {
		return nil, err
	}

	return &ast.CmdStatement{
		Command: cmd,
	}, nil
}
