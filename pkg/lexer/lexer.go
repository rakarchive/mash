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

package lexer

import (
	"unicode/utf8"

	"github.com/raklaptudirm/mash/pkg/token"
)

type lexer struct {
	src string
	ch  rune

	Tokens chan token.Token

	err ErrorHandler

	offset   int
	rdOffset int

	start token.Position
	pos   token.Position

	ErrorCount int
}

const (
	eof = -1     // end of file
	bom = 0xFEFF // byte order mark
)

type ErrorHandler func(token.Position, string)

func Lex(src string) chan token.Token {
	l := &lexer{
		src: src,

		Tokens: make(chan token.Token),

		pos: token.Position{
			Line: 1,
			Col:  1,
		},
	}
	l.run()

	return l.Tokens
}

func (l *lexer) peek() rune {
	if l.atEnd() {
		return eof
	}

	return rune(l.src[l.rdOffset])
}

func (l *lexer) consume() {
	if l.atEnd() {
		l.ch = eof
		return
	}

	r, w := rune(l.src[l.rdOffset]), 1
	if r == 0 {
		l.error("illegal character NUL")
		goto advance
	}

	if r < utf8.RuneSelf {
		goto advance
	}

	r, w = utf8.DecodeRuneInString(l.src[l.rdOffset:])

	if r == utf8.RuneError && w == 1 {
		l.error("illegal UTF-8 encoding")
		goto advance
	}

	if r == bom && l.offset > 0 {
		l.error("illegal byte order mark")
	}

advance:
	l.ch = r

	l.rdOffset += w
	l.pos.Col += w

	if r == '\n' {
		l.pos.NextLine()
	}
}

func (l *lexer) literal() string {
	return l.src[l.offset:l.rdOffset]
}

func (l *lexer) ignore() {
	l.offset = l.rdOffset
	l.start = l.pos
}

func (l *lexer) backup() {
	l.rdOffset = l.offset
}

func (l *lexer) emit(t token.TokenType) {
	l.Tokens <- token.Token{
		Type:     t,
		Literal:  l.literal(),
		Position: l.start,
	}

	l.ignore()
}

func (l *lexer) error(err string) {
	l.ErrorCount++
	if l.err != nil {
		l.err(l.pos, err)
	}
}

func (l *lexer) atEnd() bool {
	return l.rdOffset >= len(l.src)
}
