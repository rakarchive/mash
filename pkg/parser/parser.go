package parser

import (
	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/lexer"
	"github.com/raklaptudirm/mash/pkg/token"
)

type parser struct {
	tokens lexer.TokenStream

	tok token.TokenType
	pos token.Position
	lit string

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
	tok := <-p.tokens

	p.tok = tok.Type
	p.pos = tok.Position
	p.lit = tok.Literal
}

func (p *parser) error(err error) {
	p.ErrorCount++
	if p.err != nil {
		p.err(p.pos, err)
	}
}
