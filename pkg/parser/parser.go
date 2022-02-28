package parser

import (
	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/lexer"
	"github.com/raklaptudirm/mash/pkg/token"
)

type parser struct {
	tokens lexer.TokenStream

	curToken token.Token

	err lexer.ErrorHandler

	ErrorCount int
}

func Parse(t lexer.TokenStream, e lexer.ErrorHandler) *ast.Program {
	p := parser{
		tokens:     t,
		err:        e,
		ErrorCount: 0,
	}

	return p.parseProgram()
}

func (p *parser) next() {
	p.curToken = <-p.tokens
}

func (p *parser) error(err error) {
	p.ErrorCount++
	if p.err != nil {
		p.err(p.curToken.Position, err)
	}
}
