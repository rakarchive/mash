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
	"laptudirm.com/x/mash/pkg/ast"
	"laptudirm.com/x/mash/pkg/lexer"
	"laptudirm.com/x/mash/pkg/token"
)

type parser struct {
	tokens lexer.TokenStream

	// next "peek" token
	pTok token.Type
	pPos token.Position
	pLit string

	// current token
	tok token.Type
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

func (p *parser) current() token.Token {
	return token.Token{
		Type:     p.tok,
		Position: p.pos,
		Literal:  p.lit,
	}
}

func (p *parser) match(tokens ...token.Type) bool {
	if p.check(tokens...) {
		p.next()
		return true
	}

	return false
}

func (p *parser) check(tokens ...token.Type) bool {
	for _, tok := range tokens {
		if tok == p.pTok {
			return true
		}
	}

	return false
}

func (p *parser) atEnd() bool {
	return p.pTok == token.Eof
}

func (p *parser) next() {
	tok := <-p.tokens

	p.tok = p.pTok
	p.pos = p.pPos
	p.lit = p.pLit

	for tok.Type == token.Comment {
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

func (p *parser) synchronize() {
	// consume error token
	p.next()

	for !p.atEnd() {
		// semicolon ends statements, so current token is statement start
		if p.tok == token.Semicolon {
			return
		}

		switch p.pTok {
		// check for tokens which start a statement
		case token.For, token.If, token.Let, token.Break, token.Continue, token.Return:
			return
		default:
			p.next()
		}
	}
}
