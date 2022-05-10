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
	"errors"
	"fmt"
	"unicode"

	"laptudirm.com/x/mash/pkg/token"
)

var ErrEOF = errors.New("unexpected EOF")

// run starts lexing the source in l and closes the lexer's token channel
// when it is done.
func (l *lexer) run() {
	l.lexBlock(eof, token.EOF)
	close(l.Tokens)
}

func (l *lexer) lexBlock(eob rune, tok token.TokenType) {
	for {
		r := l.peek()
		switch {
		case r == eob:
			l.consume()
			l.emit(tok)
			return // block lexed

		// ignore all space runes
		case unicode.IsSpace(r):
			l.consumeAllSpace()

		case r == '#':
			l.lexComment()

		// command or statement
		default:
			if isAlphabet(r) {
				l.consumeWord()

				word := l.literal()
				// statement starts with keyword
				if token.IsKeyword(word) {
					t := token.Lookup(word)
					l.emit(t)
					l.insertSemi = t.InsertSemi()

					l.lexStmt(eob)

					// semicolon insertion
					l.emit(token.SEMICOLON)
					break
				}

				// commands don't start with a keyword
				// TODO: cleanup
				l.rdOffset = l.offset
				l.pos = l.start
			}

			l.lexCmd(eob)

			// semicolon insertion
			l.emit(token.SEMICOLON)
		}
	}
}

func isAlphabet(r rune) bool {
	return r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
}

func (l *lexer) lexStmt(eos rune) {
	for {
		l.consume()

		switch {
		case l.ch == eos:
			l.backup()
			return // will be handled by caller

		case l.ch == '\n':
			if l.insertSemi {
				// if semicolon is inserted the statement ends
				return
			}

			// ignore all space
			l.consumeAllSpace()

		case unicode.IsSpace(l.ch):
			if l.insertSemi {
				// semicolon to be inserted
				// do not ignore newlines
				l.consumeSpace()
				break
			}

			// no semicolon to be inserted
			// ignore all space runes
			l.consumeAllSpace()

		case isIdentStart(l.ch):
			t := l.lexIdent()
			l.insertSemi = t.InsertSemi()

		case l.ch == '{':
			l.emit(token.LBRACE)
			l.lexBlock('}', token.RBRACE)
			l.insertSemi = true

		case unicode.IsDigit(l.ch):
			l.lexNum()
			// semicolon should be inserted after a number
			l.insertSemi = true

		case l.ch == '"' || l.ch == '\'' || l.ch == '`':
			l.lexString()
			// semicolon should be inserted after a string
			l.insertSemi = true

		// all operator starting runes are themselves operators
		case token.IsOperator(string(l.ch)):
			t := l.lexStmtOp()
			l.insertSemi = t.InsertSemi()

		case l.ch == '#':
			// line comment
			l.lexComment()

		default:
			// rune not supported inside statements
			l.emit(token.ILLEGAL)
		}
	}
}

func (l *lexer) lexIdent() token.TokenType {
	for isIdent(l.peek()) {
		l.consume()
	}

	// lookup the token type of literal
	t := token.Lookup(l.literal())
	l.emit(t)
	return t
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func (l *lexer) lexNum() {
	base := 10 // number base

	// 0b, 0o, or 0x base specs
	if l.ch == '0' {
		var ok bool
		if base, ok = baseOf(l.peek()); ok {
			l.consume()
		}

		if l.peek() == '_' {
			l.consume()
		}
	}

	l.lexDigits(base)

	if base <= 8 {
		goto tokenize
	}

	if l.peek() == '.' {
		l.consume()
		l.lexDigits(base)
	}

	if isExponent(l.peek(), base) {
		l.consume()

		switch l.peek() {
		case '+', '-':
			l.consume()
		}

		l.lexDigits(10)
	}

tokenize:
	l.emit(token.FLOAT)
}

func baseOf(r rune) (int, bool) {
	switch r {
	case 'b', 'B':
		return 2, true
	case 'o', 'O':
		return 8, true
	case 'x', 'X':
		return 16, true
	default:
		return 8, false
	}
}

func isExponent(r rune, base int) bool {
	switch r {
	case 'p', 'P':
		return base == 16
	case 'e', 'E':
		return base == 10
	default:
		return false
	}
}

func (l *lexer) lexDigits(base int) {
	if !isBaseDigit(l.peek(), base) {
		l.error(fmt.Errorf("invalid number literal"))
		return
	}

	l.consume()

	for {
		if l.peek() == '_' {
			l.consume()
		}

		if !isBaseDigit(l.peek(), base) {
			return
		}

		l.consume()
	}
}

func isBaseDigit(r rune, b int) bool {
	switch r {
	case 'A', 'B', 'C', 'D', 'E', 'F', 'a', 'b', 'c', 'd', 'e', 'f':
		return b == 16
	case '8', '9':
		return b >= 10
	case '2', '3', '4', '5', '6', '7':
		return b >= 8
	case '0', '1':
		return true
	default:
		return false
	}
}

func (l *lexer) makeOp(target rune, pass token.TokenType, fail token.TokenType) token.TokenType {
	if l.peek() == target {
		l.consume()
		return pass
	}

	return fail
}

func (l *lexer) lexStmtOp() token.TokenType {
	var t token.TokenType
	switch l.ch {
	case '+':
		t = l.makeOp('=', token.ADD_ASSIGN, token.ADD)
	case '-':
		t = l.makeOp('=', token.SUB_ASSIGN, token.SUB)
	case '*':
		t = l.makeOp('=', token.MUL_ASSIGN, token.MUL)
	case '/':
		t = l.makeOp('=', token.QUO_ASSIGN, token.QUO)
	case '%':
		t = l.makeOp('=', token.REM_ASSIGN, token.REM)
	case '&':
		t = l.makeOp('&', token.LAND, token.AND)

		if t == token.LAND {
			break
		}

		t = l.makeOp('^', token.AND_NOT, token.AND)

		e := token.AND_NOT_ASSIGN
		if t == token.AND {
			e = token.AND_ASSIGN
		}

		t = l.makeOp('=', e, t)
	case '|':
		t = l.makeOp('|', token.LOR, token.OR)

		if t == token.OR {
			t = l.makeOp('=', token.OR_ASSIGN, token.OR)
		}
	case '^':
		t = l.makeOp('=', token.XOR_ASSIGN, token.XOR)
	case '<':
		t = l.makeOp('<', token.SHL, token.LSS)

		e := token.SHL_ASSIGN
		if t == token.LSS {
			e = token.LEQ
		}

		t = l.makeOp('=', e, t)
	case '>':
		t = l.makeOp('>', token.SHR, token.GTR)

		e := token.SHR_ASSIGN
		if t == token.GTR {
			e = token.GEQ
		}

		t = l.makeOp('=', e, t)
	case '=':
		t = l.makeOp('=', token.EQL, token.ASSIGN)
	case '!':
		t = l.makeOp('=', token.NEQ, token.NOT)
	case '(':
		t = token.LPAREN
	case '[':
		t = token.LBRACK
	case ',':
		t = token.COMMA
	case ')':
		t = token.RPAREN
	case ']':
		t = token.RBRACK
	case ';':
		t = token.SEMICOLON
	case ':':
		t = l.makeOp('=', token.DEFINE, token.COLON)
	}

	l.emit(t)
	return t
}

func (l *lexer) lexCmd(eoc rune) {
	for {
		l.consume()

		switch {
		case l.ch == eoc:
			l.backup()
			return // will be handled by lexBlock

		case l.ch == '\n':
			return // insertion in handled by lexBlock

		case unicode.IsSpace(l.ch):
			// ignore all space
			l.consumeSpace()

		case l.ch == '"':
			l.lexString()

		case l.ch == '#':
			l.lexComment()

		case isCmdOp(l.ch):
			l.lexCmdOp()

		default:
			l.consumeWord()
			l.emit(token.STRING)
		}
	}
}

func isCmdOp(r rune) bool {
	switch r {
	case '|', '&', '!':
		return true
	default:
		return false
	}
}

func (l *lexer) lexCmdOp() token.TokenType {
	var t token.TokenType
	switch l.ch {
	case '|':
		t = l.makeOp('|', token.LOR, token.OR)
	case '&':
		t = l.makeOp('&', token.LAND, token.AND)
	case '!':
		t = token.NOT
	default:
		// unreachable
		t = token.ILLEGAL
	}

	l.emit(t)
	return t
}

func (l *lexer) lexComment() {
	// consume tokens till newline or eof
	for r := l.peek(); r != '\n' && r != eof; r = l.peek() {
		l.consume()
	}

	l.emit(token.COMMENT)
}

func (l *lexer) lexString() {
	switch l.ch {
	case '`':
		l.lexRawString()
	case '\'':
		l.lexEmbeddedString()
	case '"':
		l.lexInterpretedString()
	}
}

func (l *lexer) lexRawString() {
	// consume tokens till '`' or eof
	for r := l.peek(); r != '`' && r != eof; r = l.peek() {
		l.consume()
	}

	if l.peek() == eof {
		l.error(ErrEOF)
		l.emit(token.ILLEGAL)
		return
	}

	l.consume() // consume the trailing '`'
	l.emit(token.STRING)
}

func (l *lexer) lexInterpretedString() {
	// consume tokens till '"' or eof
	for r := l.peek(); r != '"' && r != eof; r = l.peek() {
		l.consume()

		// consume escape rune
		if r == '\\' {
			l.lexStringEscape('"')
		}
	}

	if l.peek() == eof {
		l.error(ErrEOF)
		l.emit(token.ILLEGAL)
		return
	}

	l.consume() // consume the trailing '"'
	l.emit(token.STRING)
}

func (l *lexer) lexEmbeddedString() {}

var (
	ErrEsc    = errors.New("invalid escape sequence")
	ErrEscEnd = errors.New("unterminated escape sequence")
)

func (l *lexer) lexStringEscape(t rune) {
	var radix, n int
	switch l.peek() {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', t:
		l.consume()
		return
	case '0', '1', '2', '3', '4', '5', '6', '7':
		radix, n = 8, 3
	case 'x':
		radix, n = 16, 2
	case 'u':
		radix, n = 16, 4
	case 'U':
		radix, n = 16, 8
	default:
		l.error(ErrEsc)
		return
	}

	l.consume()

	for i := 0; i < n; i++ {
		r := l.peek()
		if r == eof || r == t {
			fmt.Println(r)
			l.error(ErrEscEnd)
			return
		}

		if !isBaseDigit(r, radix) {
			l.error(fmt.Errorf("illegal rune %#U in escape sequence", r))
			return
		}

		l.consume()
	}
}

// consumeSpace consumes all non newline space runes.
func (l *lexer) consumeSpace() {
	for r := l.peek(); unicode.IsSpace(r) && r != '\n'; r = l.peek() {
		l.consume()
	}

	l.ignore()
}

// consumeAllSpace consumes all space runes.
func (l *lexer) consumeAllSpace() {
	for unicode.IsSpace(l.peek()) {
		l.consume()
	}

	l.ignore()
}

func isIdent(r rune) bool {
	return isIdentStart(r) || unicode.IsDigit(r)
}

func (l *lexer) consumeWord() {
	for r := l.peek(); !unicode.IsSpace(r) && r != eof; r = l.peek() {
		l.consume()
	}
}
