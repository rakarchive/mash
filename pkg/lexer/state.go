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
	"unicode"

	"github.com/raklaptudirm/mash/pkg/token"
)

func (l *lexer) run() {
	lexBase(l)
	close(l.Tokens)
}

func lexBase(l *lexer) {
next:
	r := l.peek()
	switch {
	case r == eof:
		l.emit(token.SEMICOLON)
		l.emit(token.EOF)
		return
	case unicode.IsSpace(r):
		consumeAllSpace(l)
	case r == '#':
		lexComment(l)
	case isAlphabet(r):
		consumeWord(l)

		word := l.literal()
		if token.IsKeyword(word) {
			l.emit(token.Lookup(word))
			lexStmt(l)
			break
		}

		l.emit(cmdOpLookup(word))
		fallthrough
	default:
		lexCmd(l)
	}

	goto next
}

func isAlphabet(r rune) bool {
	return r > 'A' && r < 'Z' || r > 'a' && r < 'z'
}

func lexStmt(l *lexer) {
next:
	l.consume()

	switch {
	case l.ch == '\n':
		// semicolon insertion
		if l.prev.InsertSemi() {
			l.emit(token.SEMICOLON)
			return
		}

		consumeAllSpace(l)
	case unicode.IsSpace(l.ch):
		// ignore whitespace
		consumeSpace(l)

	// literals
	case isIdentStart(l.ch):
		// identifier
		lexIdent(l)
	case unicode.IsDigit(l.ch):
		// number
		lexNum(l)
	case l.ch == '"':
		// format string
		lexString(l)

	// operators
	case token.IsOperator(string(l.ch)):
		lexStmtOp(l)

	// special
	case l.ch == '#':
		// line comment
		lexComment(l)
	case l.ch == eof:
		return
	default:
		// rune not supported
		l.emit(token.ILLEGAL)
	}

	goto next
}

func lexIdent(l *lexer) {
	for isIdent(l.peek()) {
		l.consume()
	}

	l.emit(token.Lookup(l.literal()))
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func lexNum(l *lexer) {
	// TODO: support more number types
	for unicode.IsDigit(l.peek()) {
		l.consume()
	}

	l.emit(token.FLOAT)
}

func (l *lexer) makeOp(target rune, pass token.TokenType, fail token.TokenType) token.TokenType {
	if l.peek() == target {
		l.consume()
		return pass
	}

	return fail
}

func lexStmtOp(l *lexer) {
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
		t = l.makeOp('=', token.ASSIGN, token.EQL)
	case '!':
		t = l.makeOp('=', token.NEQ, token.NOT)
	case '(':
		t = token.LPAREN
	case '[':
		t = token.LPAREN
	case '{':
		t = token.LBRACE
	case ',':
		t = token.COMMA
	case ')':
		t = token.RPAREN
	case ']':
		t = token.RBRACK
	case '}':
		t = token.RBRACE
	case ';':
		t = token.SEMICOLON
	case ':':
		t = l.makeOp('=', token.DEFINE, token.COLON)
	}

	l.emit(t)
}

func lexCmd(l *lexer) {
next:
	l.consume()

	switch {
	// semicolon insertion
	case l.ch == '\n':
		l.emit(token.SEMICOLON)
		return

	case unicode.IsSpace(l.ch):
		// ignore whitespace
		consumeSpace(l)

	case l.ch == '"':
		lexString(l)

	// comment
	case l.ch == '#':
		lexComment(l)
	case l.ch == eof:
		return
	default:
		consumeWord(l)
		l.emit(cmdOpLookup(l.literal()))
	}

	goto next
}

func cmdOpLookup(s string) token.TokenType {
	switch s {
	case "||":
		return token.LOR
	case "&&":
		return token.LAND
	case "!":
		return token.NOT
	case "<":
		return token.LSS
	case ">":
		return token.GTR
	case ">>":
		return token.SHR
	case "|":
		return token.OR
	default:
		return token.STRING
	}
}

func lexComment(l *lexer) {
	for r := l.peek(); r != '\n' && r != eof; r = l.peek() {
		l.consume()

		if r == '\\' {
			if l.peek() == eof {
				l.error("unexpected EOF")
				break
			}

			l.consume()
		}
	}

	l.emit(token.COMMENT)
}

func lexString(l *lexer) {
	for r := l.peek(); r != '"' && r != eof; r = l.peek() {
		l.consume()

		if r == '\\' {
			if l.peek() == eof {
				l.error("unexpected EOF")
				break
			}

			l.consume()
		}
	}

	if l.peek() == eof {
		l.error("unexpected EOF")
	}

	l.emit(token.STRING)
}

func consumeSpace(l *lexer) {
	for r := l.peek(); unicode.IsSpace(r) && r != '\n'; r = l.peek() {
		l.consume()
	}

	l.ignore()
}

func consumeAllSpace(l *lexer) {
	for unicode.IsSpace(l.peek()) {
		l.consume()
	}

	l.ignore()
}

func isIdent(r rune) bool {
	return isIdentStart(r) || unicode.IsDigit(r)
}

func consumeWord(l *lexer) {
	for r := l.peek(); !unicode.IsSpace(r) && r != eof; r = l.peek() {
		l.consume()
	}
}
