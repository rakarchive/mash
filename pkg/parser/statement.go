package parser

import (
	"fmt"

	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/token"
)

// program = { statement } EOF
func (p *parser) parseProgram() (*ast.Program, error) {
	program := &ast.Program{}

	for !p.atEnd() {
		stmt, err := p.parseStatement()
		if err != nil {
			p.error(p.pPos, err)
			// sync parser to avoid cascading errors
			p.synchronize()
			continue
		}

		program.Statements = append(program.Statements, stmt)
	}
	return program, nil
}

// statement = ( blockStmt | letStmt | ifStmt | cmdStmt ) ";"
func (p *parser) parseStatement() (ast.Statement, error) {
	var stmt ast.Statement
	var err error

	switch p.pTok {
	case token.LBRACE:
		stmt, err = p.parseBlockStmt()
	case token.LET:
		stmt, err = p.parseLetStmt()
	case token.IF:
		stmt, err = p.parseIfStatement()
	case token.FOR:
		stmt, err = p.parseForStmt()
	case token.STRING, token.NOT:
		stmt, err = p.parseCmdStmt()
	default:
		return nil, fmt.Errorf("illegal token %s at line start", p.pTok)
	}

	// only check for semicolons if no errors have occurred
	if err == nil && !p.match(token.SEMICOLON) {
		return nil, fmt.Errorf("expected ';', received %s", p.pTok)
	}
	return stmt, err
}

// blockStmt = "{" { statement } "}"
func (p *parser) parseBlockStmt() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{}

	if !p.match(token.LBRACE) {
		return nil, fmt.Errorf("expected '{', received %s", p.pTok)
	}

	for p.pTok != token.RBRACE && !p.atEnd() {
		stmt, err := p.parseStatement()
		if err != nil {
			p.error(p.pPos, err)
			// sync parser to avoid cascading errors
			p.synchronize()
			continue
		}

		block.Statements = append(block.Statements, stmt)
	}

	if !p.match(token.RBRACE) {
		return nil, fmt.Errorf("expected '}', received %s", p.pTok)
	}
	return block, nil
}

// letStmt = "let" assignExpr
func (p *parser) parseLetStmt() (*ast.LetStatement, error) {
	if !p.match(token.LET) {
		return nil, fmt.Errorf("expected 'let', received %s", p.pTok)
	}

	expr, err := p.parseExprAssign()
	if err != nil {
		return nil, err
	}

	return &ast.LetStatement{
		Expression: expr,
	}, nil
}

// ifStmt = "if" expression blockStmt { "elif" expression blockStmt } [ "else" blockStmt ]
func (p *parser) parseIfStatement() (*ast.IfStatement, error) {
	if !p.match(token.IF) {
		return nil, fmt.Errorf("expected 'if', received %s", p.pTok)
	}

	cond, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlockStmt()
	if err != nil {
		return nil, err
	}

	stmt := ast.IfStatement{
		Condition: cond,
		BlockStmt: block,
	}

	for p.match(token.ELIF) {
		cond, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		block, err := p.parseBlockStmt()
		if err != nil {
			return nil, err
		}

		stmt.ElifBlock = append(stmt.ElifBlock, ast.ElifBlock{
			Condition: cond,
			BlockStmt: block,
		})
	}

	if p.match(token.ELSE) {
		block, err := p.parseBlockStmt()
		if err != nil {
			return nil, err
		}

		stmt.ElseBlock = block
	}

	return &stmt, nil
}

// forStmt = "for" [ expression ] blockStmt
func (p *parser) parseForStmt() (*ast.ForStatement, error) {
	if !p.match(token.FOR) {
		return nil, fmt.Errorf("expected 'for', received %s", p.pTok)
	}

	var condition ast.Expression
	var err error

	if !p.check(token.LBRACE) {
		condition, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	block, err := p.parseBlockStmt()
	if err != nil {
		return nil, err
	}

	return &ast.ForStatement{
		Condition: condition,
		BlockStmt: block,
	}, nil
}

func (p *parser) parseCmdStmt() (*ast.CmdStatement, error) {
	cmd, err := p.parseCommand()
	if err != nil {
		return nil, err
	}

	return &ast.CmdStatement{
		Command: cmd,
	}, nil
}
