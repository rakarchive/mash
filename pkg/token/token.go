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

package token

import (
	"strconv"
	"unicode"
)

// TokenType represents the type of a token which will be emitted by the
// lexer.
type TokenType int

// Various types of tokens emitted by the lexer.
const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	IDENT  // main
	FLOAT  // 3.14
	STRING // "abc"
	literal_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND // &&
	LOR  // ||

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	DEFINE // :=
	NOT    // !

	NEQ // !=
	LEQ // <=
	GEQ // >=

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	operator_end

	keyword_beg
	// Keywords
	FOR
	IF
	ELIF
	ELSE

	LET
	FUNC

	BREAK
	CONTINUE
	RETURN
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND:     "&",
	OR:      "|",
	XOR:     "^",
	SHL:     "<<",
	SHR:     ">>",
	AND_NOT: "&^",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	AND_ASSIGN:     "&=",
	OR_ASSIGN:      "|=",
	XOR_ASSIGN:     "^=",
	SHL_ASSIGN:     "<<=",
	SHR_ASSIGN:     ">>=",
	AND_NOT_ASSIGN: "&^=",

	LAND: "&&",
	LOR:  "||",

	EQL:    "==",
	LSS:    "<",
	GTR:    ">",
	ASSIGN: "=",
	NOT:    "!",

	NEQ:    "!=",
	LEQ:    "<=",
	GEQ:    ">=",
	DEFINE: ":=",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	FOR:  "for",
	IF:   "if",
	ELIF: "elif",
	ELSE: "else",

	LET:  "let",
	FUNC: "func",

	BREAK:    "break",
	CONTINUE: "continue",
	RETURN:   "return",
}

func token(s string) TokenType {
	for t, val := range tokens {
		if val == s {
			return TokenType(t)
		}
	}

	return ILLEGAL
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok TokenType) String() string {
	s := ""
	if 0 <= tok && tok < TokenType(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// InsertSemi returns a boolean depending on wether a semicolon
// should be inserted after a token of type tok. It returns true if
// a semicolon should be inserted, and false if should not.
//
func (tok TokenType) InsertSemi() bool {
	if tok.IsLiteral() {
		return true
	}

	switch tok {
	case RPAREN, RBRACK, RBRACE, BREAK, CONTINUE, RETURN:
		return true
	default:
		return false
	}
}

// IsLiteral returns a boolean depending on wether the type of tok is
// a valid literal. Literals are tokens of with a value greater than
// literal_beg but less than literal_end.
//
func (tok TokenType) IsLiteral() bool {
	return literal_beg < tok && tok < literal_end
}

// IsOperator returns a boolean depending on wether the type of tok is
// a valid operator. Operators are tokens of with a value greater than
// operator_beg but less than operator_end.
//
func (tok TokenType) IsOperator() bool {
	return operator_beg < tok && tok < operator_end
}

// IsKeyword returns a boolean depending on wether the type of tok is
// a valid keyword. Keywords are tokens of with a value greater than
// keyword_beg but less than keyword_end.
//
func (tok TokenType) IsKeyword() bool {
	return keyword_beg < tok && tok < keyword_end
}

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// IsKeyword returns a boolean depending on wether name is a valid
// keyword. A string is a keyword if it is present in the keywords
// map.
func IsKeyword(name string) bool {
	_, ok := keywords[name]
	return ok
}

// IsIdentifier returns a boolean depending of wether name is a valid
// identifier. A string is a valid identifier if it's first letter is
// an unicode letter(gc = L) or an underscore, while the rest of the
// characters are letters, underscores, or unicode decimal digits
// (gc =  Nd).
func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && (i == 0 || !unicode.IsDigit(c)) && c != '_' {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}

// IsOperator returns a boolean depending on wether name is a valid
// operator or not. If the string belongs in the list of mash operators,
// it is a valid operator.
func IsOperator(s string) bool {
	t := token(s)
	return t.IsOperator()
}

// Lookup checks if name is a keyword, and returns the token type of the
// keyword if it is. Otherwise, it returns IDENT.
func Lookup(name string) TokenType {
	if tok, ok := keywords[name]; ok {
		return tok
	}

	return IDENT
}

// Token represtents a single token which will be emitted by the lexer.
type Token struct {
	Type     TokenType // type of the token
	Literal  string    // literal in source
	Position Position  // position in source
}
