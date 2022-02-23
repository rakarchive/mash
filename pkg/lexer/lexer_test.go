package lexer_test

import (
	"testing"

	"github.com/raklaptudirm/mash/pkg/lexer"
	"github.com/raklaptudirm/mash/pkg/token"
)

func TestLexer(t *testing.T) {
	input := `# comment line
for
if
elif
else

let
func

break
continue
return

let identifier
let 3141592653
let "a string"

# line comment
let i # inline

let +
let -
let *
let /
let %

let &
let |
let ^
let <<
let >>
let &^

let +=
let -=
let *=
let /=
let %=

let &=
let |=
let ^=
let <<=
let >>=
let &^=

let &&
let ||

let ==
let <
let >
let =
let :=
let !

let !=
let <=
let >=

let (
let [
let {
let ,

let )
let ]
let }
let ;
let :

break

echo a command
||&&!>>>|
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedCol     int
	}{
		{token.COMMENT, "# comment line", 1, 1},
		{token.FOR, "for", 2, 1},
		{token.IF, "if", 3, 1},
		{token.ELIF, "elif", 4, 1},
		{token.ELSE, "else", 5, 1},
		{token.LET, "let", 7, 1},
		{token.FUNC, "func", 8, 1},
		{token.BREAK, "break", 10, 1},
		{token.SEMICOLON, "\n", 10, 6},
		{token.CONTINUE, "continue", 11, 1},
		{token.SEMICOLON, "\n", 11, 9},
		{token.RETURN, "return", 12, 1},
		{token.SEMICOLON, "\n", 12, 7},
		{token.LET, "let", 14, 1},
		{token.IDENT, "identifier", 14, 5},
		{token.SEMICOLON, "\n", 14, 15},
		{token.LET, "let", 15, 1},
		{token.FLOAT, "3141592653", 15, 5},
		{token.SEMICOLON, "\n", 15, 15},
		{token.LET, "let", 16, 1},
		{token.STRING, "\"a string\"", 16, 5},
		{token.SEMICOLON, "\n", 16, 15},
		{token.COMMENT, "# line comment", 18, 1},
		{token.LET, "let", 19, 1},
		{token.IDENT, "i", 19, 5},
		{token.COMMENT, "# inline", 19, 7},
		{token.SEMICOLON, "\n", 19, 15},
		{token.LET, "let", 21, 1},
		{token.ADD, "+", 21, 5},
		{token.LET, "let", 22, 1},
		{token.SUB, "-", 22, 5},
		{token.LET, "let", 23, 1},
		{token.MUL, "*", 23, 5},
		{token.LET, "let", 24, 1},
		{token.QUO, "/", 24, 5},
		{token.LET, "let", 25, 1},
		{token.REM, "%", 25, 5},
		{token.LET, "let", 27, 1},
		{token.AND, "&", 27, 5},
		{token.LET, "let", 28, 1},
		{token.OR, "|", 28, 5},
		{token.LET, "let", 29, 1},
		{token.XOR, "^", 29, 5},
		{token.LET, "let", 30, 1},
		{token.SHL, "<<", 30, 5},
		{token.LET, "let", 31, 1},
		{token.SHR, ">>", 31, 5},
		{token.LET, "let", 32, 1},
		{token.AND_NOT, "&^", 32, 5},
		{token.LET, "let", 34, 1},
		{token.ADD_ASSIGN, "+=", 34, 5},
		{token.LET, "let", 35, 1},
		{token.SUB_ASSIGN, "-=", 35, 5},
		{token.LET, "let", 36, 1},
		{token.MUL_ASSIGN, "*=", 36, 5},
		{token.LET, "let", 37, 1},
		{token.QUO_ASSIGN, "/=", 37, 5},
		{token.LET, "let", 38, 1},
		{token.REM_ASSIGN, "%=", 38, 5},
		{token.LET, "let", 40, 1},
		{token.AND_ASSIGN, "&=", 40, 5},
		{token.LET, "let", 41, 1},
		{token.OR_ASSIGN, "|=", 41, 5},
		{token.LET, "let", 42, 1},
		{token.XOR_ASSIGN, "^=", 42, 5},
		{token.LET, "let", 43, 1},
		{token.SHL_ASSIGN, "<<=", 43, 5},
		{token.LET, "let", 44, 1},
		{token.SHR_ASSIGN, ">>=", 44, 5},
		{token.LET, "let", 45, 1},
		{token.AND_NOT_ASSIGN, "&^=", 45, 5},
		{token.LET, "let", 47, 1},
		{token.LAND, "&&", 47, 5},
		{token.LET, "let", 48, 1},
		{token.LOR, "||", 48, 5},
		{token.LET, "let", 50, 1},
		{token.EQL, "==", 50, 5},
		{token.LET, "let", 51, 1},
		{token.LSS, "<", 51, 5},
		{token.LET, "let", 52, 1},
		{token.GTR, ">", 52, 5},
		{token.LET, "let", 53, 1},
		{token.ASSIGN, "=", 53, 5},
		{token.LET, "let", 54, 1},
		{token.DEFINE, ":=", 54, 5},
		{token.LET, "let", 55, 1},
		{token.NOT, "!", 55, 5},
		{token.LET, "let", 57, 1},
		{token.NEQ, "!=", 57, 5},
		{token.LET, "let", 58, 1},
		{token.LEQ, "<=", 58, 5},
		{token.LET, "let", 59, 1},
		{token.GEQ, ">=", 59, 5},
		{token.LET, "let", 61, 1},
		{token.LPAREN, "(", 61, 5},
		{token.LET, "let", 62, 1},
		{token.LBRACK, "[", 62, 5},
		{token.LET, "let", 63, 1},
		{token.LBRACE, "{", 63, 5},
		{token.LET, "let", 64, 1},
		{token.COMMA, ",", 64, 5},
		{token.LET, "let", 66, 1},
		{token.RPAREN, ")", 66, 5},
		{token.SEMICOLON, "\n", 66, 6},
		{token.LET, "let", 67, 1},
		{token.RBRACK, "]", 67, 5},
		{token.SEMICOLON, "\n", 67, 6},
		{token.LET, "let", 68, 1},
		{token.SEMICOLON, "", 68, 5},
		{token.RBRACE, "}", 68, 5},
		{token.SEMICOLON, "\n", 68, 6},
		{token.LET, "let", 69, 1},
		{token.SEMICOLON, ";", 69, 5},
		{token.LET, "let", 70, 1},
		{token.COLON, ":", 70, 5},
		{token.BREAK, "break", 72, 1},
		{token.SEMICOLON, "\n", 72, 6},
		{token.STRING, "echo", 74, 1},
		{token.STRING, "a", 74, 6},
		{token.STRING, "command", 74, 8},
		{token.SEMICOLON, "\n", 74, 15},
		{token.LOR, "||", 75, 1},
		{token.LAND, "&&", 75, 3},
		{token.NOT, "!", 75, 5},
		{token.SHR, ">>", 75, 6},
		{token.GTR, ">", 75, 8},
		{token.OR, "|", 75, 9},
		{token.SEMICOLON, "\n", 75, 10},
		{token.EOF, "", 76, 1},
	}

	index := 0
	for token := range lexer.Lex(input, nil) {
		t.Logf("%s %s %v\n", &token.Position, token.Type, token.Literal)
		if token.Type != tests[index].expectedType {
			t.Fatalf("case %v: expected token type %q, got %q", index, tests[index].expectedType, token.Type)
		}
		if token.Literal != tests[index].expectedLiteral {
			t.Fatalf("case %v: expected token literal %q, got %q", index, tests[index].expectedLiteral, token.Literal)
		}
		if token.Position.Line != tests[index].expectedLine {
			t.Fatalf("case %v: expected token line %d, got %d", index, tests[index].expectedLine, token.Position.Line)
		}
		if token.Position.Col != tests[index].expectedCol {
			t.Fatalf("case %v: expected token col %d, got %d", index, tests[index].expectedCol, token.Position.Col)
		}
		index++
	}

}
