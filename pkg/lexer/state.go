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

type stateFunc func(*lexer) stateFunc

func (l *lexer) run() {
	for state := lexBase; state != nil; {
		state = state(l)
	}
	close(l.Tokens)
}

func lexBase(l *lexer) stateFunc {
	r := l.peek()
	if unicode.IsSpace(r) {
		l.consumeSpace()
	}

	if isAlphabet(r) {
		l.consumeWord()

		word := l.literal()
		if token.IsKeyword(word) {
			l.emit(token.Lookup(word))
			return lexStmt
		}

		l.backup()
	}

	return lexCmd
}

func lexStmt(l *lexer) stateFunc {
	l.consume()

	switch {
	case l.ch == '\n':
		// semicolon insertion
		if l.prev.InsertSemi() {
			l.emit(token.SEMICOLON)
			return lexBase
		}

		l.consumeSpace()

	case unicode.IsSpace(l.ch):
		// ignore whitespace
		l.consumeSpace()

	// literals
	case isIdentStart(l.ch):
		// identifier
		l.consumeIdent()
		l.emit(token.Lookup(l.literal()))
	case unicode.IsDigit(l.ch):
		// number
		return lexNum
	case l.ch == '"':
		// format string
		l.consumeString()
		l.emit(token.STRING)

	// operators
	case token.IsOperator(string(l.ch)):
		return lexStmtOp

	// special
	case l.ch == '#':
		// line comment
		l.consumeComment()
		l.emit(token.COMMENT)
	case l.ch == eof:
		l.emit(token.EOF)
		return nil
	default:
		// rune not supported
		l.emit(token.ILLEGAL)
	}

	return lexStmt
}

func lexNum(l *lexer) stateFunc {
	// TODO: support more number types
	for unicode.IsDigit(l.peek()) {
		l.consume()
	}

	l.emit(token.FLOAT)
	return lexStmt
}

func (l *lexer) makeOp(target rune, pass token.TokenType, fail token.TokenType) token.TokenType {
	if l.peek() == target {
		l.consume()
		return pass
	}

	return fail
}

func lexStmtOp(l *lexer) stateFunc {
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
	return lexStmt
}

func lexCmd(l *lexer) stateFunc {
	l.consume()

	switch {
	// semicolon insertion
	case l.ch == '\n':
		l.emit(token.SEMICOLON)
		return lexBase

	case unicode.IsSpace(l.ch):
		// ignore whitespace
		l.consumeSpace()

	case l.ch == '"':
		l.consumeString()
		l.emit(token.STRING)

	// comment
	case l.ch == '#':
		l.consumeComment()
		l.emit(token.COMMENT)

	default:
		l.consumeWord()
		l.emit(cmdOpLookup(l.literal()))
	}

	return lexCmd
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

func isAlphabet(r rune) bool {
	return r > 'A' && r < 'Z' || r > 'a' && r < 'z'
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdent(r rune) bool {
	return isIdentStart(r) || unicode.IsDigit(r)
}
