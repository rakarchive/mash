package lexer_test

import (
	"testing"

	"github.com/raklaptudirm/mash/pkg/lexer"
)

func TestLexerSimpleInputs(t *testing.T) {
	tests := []struct {
		input         string
		expectedType  lexer.TokenType
		expectedValue string
	}{
		{";", lexer.SEMICOLON, ";"},
		{">", lexer.GREATER, ">"},
		{">>", lexer.GREATGREAT, ">>"},
		{"<", lexer.LESS, "<"},
		{">&", lexer.GREATAMPERSAND, ">&"},
		{"+", lexer.ILLEGAL, "+"},
		{";", lexer.SEMICOLON, ";"},
		{">", lexer.GREATER, ">"},
		{"<", lexer.LESS, "<"},
		{">>", lexer.GREATGREAT, ">>"},
		{">&", lexer.GREATAMPERSAND, ">&"},
		{"<&", lexer.LESSAMPERSAND, "<&"},
		{"|", lexer.PIPE, "|"},
		{"&", lexer.AMPERSAND, "&"},
		{"haha", lexer.IDENT, "haha"},
		{"`", lexer.ILLEGAL, "`"},
		{"'", lexer.ILLEGAL, "'"},
		{"\"", lexer.ILLEGAL, "\""},
		{"# \n", lexer.COMMENT, "# \n"},
		{"`haha`", lexer.BACKQUOTE, "`haha`"},
		{"'haha'", lexer.SINGLEQUOTE, "'haha'"},
		{"\"haha\"", lexer.DOUBLEQUOTE, "\"haha\""},
	}
	for _, test := range tests {
		l := lexer.Lex(test.input)
		for c := range l.Tokens {
			if c.Type != test.expectedType {
				t.Errorf("Expected type %v, got %v", test.expectedType, c.Type)
			}
			if c.Val != test.expectedValue {
				t.Errorf("Expected value %q, got %q", test.expectedValue, c.Val)
			}
		}
	}
}

func TestLexerMultiTokenInput(t *testing.T) {
	input := `; > < >> >& <& | & haha # 
;  >   >> "something" 'haha'` + " `blah blah` "
	tests := []struct {
		expectedType  lexer.TokenType
		expectedValue string
		expectedLine  int
		expectedPos   int
		expectedCol   int
	}{
		{lexer.SEMICOLON, ";", 0, 1, 1},
		{lexer.GREATER, ">", 0, 3, 3},
		{lexer.LESS, "<", 0, 5, 5},
		{lexer.GREATGREAT, ">>", 0, 8, 8},
		{lexer.GREATAMPERSAND, ">&", 0, 11, 11},
		{lexer.LESSAMPERSAND, "<&", 0, 14, 14},
		{lexer.PIPE, "|", 0, 16, 16},
		{lexer.AMPERSAND, "&", 0, 18, 18},
		{lexer.IDENT, "haha", 0, 23, 23},
		{lexer.COMMENT, "# \n", 1, 27, 0},
		{lexer.SEMICOLON, ";", 1, 28, 1},
		{lexer.GREATER, ">", 1, 31, 4},
		{lexer.GREATGREAT, ">>", 1, 36, 9},
		{lexer.DOUBLEQUOTE, "\"something\"", 1, 48, 21},
		{lexer.SINGLEQUOTE, "'haha'", 1, 55, 28},
		{lexer.BACKQUOTE, "`blah blah`", 1, 67, 40},
	}
	l := lexer.Lex(input)
	index := 0
	for c := range l.Tokens {
		if c.Type != tests[index].expectedType {
			t.Errorf("Expected type %q, got %q at index %v", tests[index].expectedType, c.Type, index)
		}
		if c.Val != tests[index].expectedValue {
			t.Errorf("Expected value %q, got %q at index %v", tests[index].expectedValue, c.Val, index)
		}
		if c.Pos != tests[index].expectedPos {
			t.Errorf("Expected pos %v, got %v at index %v", tests[index].expectedPos, c.Pos, index)
		}
		if c.Line != tests[index].expectedLine {
			t.Errorf("Expected line %v, got %v at index %v", tests[index].expectedLine, c.Line, index)
		}
		if c.Col != tests[index].expectedCol {
			t.Errorf("Expected col %v, got %v at index %v", tests[index].expectedCol, c.Col, index)
		}
		index++
	}
}
