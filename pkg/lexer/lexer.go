package lexer

type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	SEMICOLON // ;
	NEWLINE   // \n

	// Identifiers and literals
	IDENT

	// Logical operators
	AND
	OR

	// Quote operators
	SINGLE_QUOTE
	DOUBLE_QUOTE

	// Redirection operators
	TO_RIGHT
	TO_LEFT

	// Pipe operators
	PIPE
)

func (t TokenType) String() string {
	var tokenMap = map[TokenType]string{
		EOF:          "EOF",
		SEMICOLON:    ";",
		NEWLINE:      "\n",
		IDENT:        "IDENT",
		AND:          "&&",
		OR:           "||",
		SINGLE_QUOTE: "'",
		DOUBLE_QUOTE: "\"",
		TO_RIGHT:     ">",
		TO_LEFT:      "<",
		PIPE:         "|",
	}

	return tokenMap[t]
}

type Token struct {
	Type  TokenType
	Value string
}

// StateFn this is a recursive definition of the State - makes it easy to navigate ast type structures.
type StateFn func(*Lexer) StateFn

type Lexer struct {
	Name   string
	Input  string
	State  StateFn
	Tokens chan Token
	Start  int
	Pos    int
	End    int
}

func (l *Lexer) Emit(tokenType TokenType) {
	l.Tokens <- Token{
		Type:  tokenType,
		Value: l.Input[l.Start:l.Pos],
	}
	l.Start = l.Pos
}

func (l *Lexer) Inc() {
	l.Pos++
	if l.Pos >= len(l.Input) {
		l.Emit(EOF)
	}
}

func (l *Lexer) InputToEnd() string {
	return l.Input[l.Pos:]
}
