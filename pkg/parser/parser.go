package parser

import (
	"fmt"

	"github.com/raklaptudirm/mash/pkg/ast"
	"github.com/raklaptudirm/mash/pkg/lexer"
	"github.com/raklaptudirm/mash/pkg/token"
)

type parser struct {
	tokens lexer.TokenStream

	pTok token.TokenType
	pPos token.Position
	pLit string

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

	// start token consumption
	p.next()

	return p.parseProgram()
}

func (p *parser) consume(t token.TokenType) {
	if !p.match(t) {
		p.error(p.pPos, fmt.Errorf("expected %s, received %s", t, p.pTok))
	}
}

func (p *parser) current() token.Token {
	return token.Token{
		Type:     p.tok,
		Position: p.pos,
		Literal:  p.lit,
	}
}

func (p *parser) match(tokens ...token.TokenType) bool {
	for _, tok := range tokens {
		if p.check(tok) {
			p.next()
			return true
		}
	}

	return false
}

func (p *parser) check(tok token.TokenType) bool {
	if p.pTok == token.EOF {
		return false
	}

	return tok == p.pTok
}

func (p *parser) atEnd() bool {
	return p.pTok == token.EOF
}

func (p *parser) next() {
	tok := <-p.tokens

	p.tok = p.pTok
	p.pos = p.pPos
	p.lit = p.pLit

	for tok.Type == token.COMMENT {
		tok = <-p.tokens
	}

	p.pTok = tok.Type
	p.pPos = tok.Position
	p.pLit = tok.Literal
}

func (p *parser) error(pos token.Position, err error) {
	p.ErrorCount++
	if p.err != nil {
		p.err(pos, err)
	}
}
