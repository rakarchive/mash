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
		{"\"", lexer.DOUBLEQUOTE, "\""},
		{"'", lexer.SINGLEQUOTE, "'"},
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
		{"`", lexer.BACKQUOTE, "`"},
		{"# \n", lexer.COMMENT, "# \n"},
	}
	for _, test := range tests {
		l := lexer.Lex(test.input)
		for c := range l.Tokens {
			if c.Type != test.expectedType {
				t.Errorf("Expected type %q, got %q", test.expectedType, c.Type)
			}
			if c.Val != test.expectedValue {
				t.Errorf("Expected value %q, got %q", test.expectedValue, c.Val)
			}
		}
	}
}

func TestLexerMultiTokenInput(t *testing.T) {
	input := `; > < >> >& <& | & haha # 
;  >   >>`
	tests := []struct {
		expectedType  lexer.TokenType
		expectedValue string
	}{
		{lexer.SEMICOLON, ";"},
		{lexer.GREATER, ">"},
		{lexer.LESS, "<"},
		{lexer.GREATGREAT, ">>"},
		{lexer.GREATAMPERSAND, ">&"},
		{lexer.LESSAMPERSAND, "<&"},
		{lexer.PIPE, "|"},
		{lexer.AMPERSAND, "&"},
		{lexer.IDENT, "haha"},
		{lexer.COMMENT, "# \n"},
		{lexer.SEMICOLON, ";"},
		{lexer.GREATER, ">"},
		{lexer.GREATGREAT, ">>"},
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
		index++
	}
}
