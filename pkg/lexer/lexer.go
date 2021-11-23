package lexer

import (
	"unicode"
	"unicode/utf8"
)

type TokenType int

// These are the main token types for shell
const (
	// Special tokens
	ILLEGAL TokenType = iota
	SEMICOLON
	GREATER
	LESS
	GREATGREAT
	LESSLESS
	GREATAMPERSAND
	LESSAMPERSAND
	PIPE
	AMPERSAND
	IDENT
	SINGLEQUOTE
	DOUBLEQUOTE
	BACKQUOTE
	COMMENT
	EOF
)

type Token struct {
	Type TokenType
	Val  string
}

// grammar
// args -> IDENTS | empty
// cmd and args -> IDENT args
// pipe list -> PIPE cmd and args | cmd and args
// io modifier -> GREAT IDENT | GREATGREAT IDENT | ...
// io modifiers list -> io modifiers | empty
// background optional -> ampersand | empty
// cmd line -> pipe list io modifiers background optional NEWLINE | NEWLINE | ILLEGAL recovery.

type Lexer struct {
	input  string // input string
	start  int    // start position of current token
	pos    int    // current position in input
	width  int    // width of the last rune read
	Tokens chan Token
}

func Lex(input string) *Lexer {
	l := &Lexer{
		input:  input,
		Tokens: make(chan Token),
	}
	go l.run()
	return l
}

func (l *Lexer) run() {
	for state := lexCommand; state != nil; {
		state = state(l)
	}
	close(l.Tokens)
}

func (l *Lexer) emit(t TokenType) {
	l.Tokens <- Token{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return -1 // EOF
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

type stateFn func(*Lexer) stateFn

func lexCommand(l *Lexer) stateFn {
	switch r := l.next(); {
	case r == ' ' || r == '\t' || r == '\r' || r == '\n':
		l.ignore()
	case r == '|':
		l.emit(PIPE)
	case r == '&':
		l.emit(AMPERSAND)
	case r == '<':
		return lexLess
	case r == ';':
		l.emit(SEMICOLON)
	case r == '>':
		return lexGreater
	case r == '\'':
		return lexSingleQuote
	case r == '"':
		return lexDoubleQuote
	case r == '`':
		return lexBackQuote
	case r == '$':
		return lexVariable
	case r == '#':
		return lexComment
	case isAlphaNumeric(r):
		l.backup()
		return lexIdent
	case r == -1:
		return nil
	default:
		return l.errorf("unexpected character %#U", r)
	}
	return lexCommand
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || r == '-' || r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.Tokens <- Token{
		Type: ILLEGAL,
		Val:  l.input[l.start:l.pos],
	}
	return nil
}

func lexComment(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == '\n' {
			break
		}
		if r == -1 {
			break
		}
	}
	l.emit(COMMENT)
	return lexCommand
}

func lexSingleQuote(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == '\'' {
			break
		}
		if r == -1 {
			return l.errorf("unexpected EOF")
		}
	}
	l.emit(SINGLEQUOTE)
	return lexCommand
}

func lexDoubleQuote(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == '"' {
			break
		}
		if r == -1 {
			return l.errorf("unexpected EOF")
		}
	}
	l.emit(DOUBLEQUOTE)
	return lexCommand
}

func lexBackQuote(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == '`' {
			break
		}
		if r == -1 {
			return l.errorf("unexpected EOF")
		}
	}
	l.emit(BACKQUOTE)
	return lexCommand
}

func lexVariable(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '|' || r == '&' || r == ';' || r == '<' || r == '>' || r == '"' || r == '\'' || r == '`' || r == '$' || r == '#' || r == -1 {
			break
		}
	}
	l.emit(IDENT)
	return lexCommand
}

func lexIdent(l *Lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			l.emit(IDENT)
			break Loop
		}
	}
	return lexCommand
}

func lexGreater(l *Lexer) stateFn {
	r := l.peek()
	switch r {
	case '&':
		l.next()
		l.emit(GREATAMPERSAND)
	case '>':
		l.next()
		l.emit(GREATGREAT)
	default:
		l.emit(GREATER)
	}
	return lexCommand
}

func lexLess(l *Lexer) stateFn {
	r := l.peek()
	switch r {
	case '&':
		l.next()
		l.emit(LESSAMPERSAND)
	case '<':
		l.next()
		l.emit(LESSLESS)
	default:
		l.emit(LESS)
	}
	return lexCommand
}
