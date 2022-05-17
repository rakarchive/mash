package ast

import "laptudirm.com/x/mash/pkg/token"

type Literal interface {
	Expression
	Literal()
}

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteral) Node()       {}
func (n *NumberLiteral) Literal()    {}
func (n *NumberLiteral) Expression() {}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (n *StringLiteral) Node()             {}
func (n *StringLiteral) Literal()          {}
func (n *StringLiteral) Expression()       {}
func (n *StringLiteral) CommandComponent() {}

type FunctionLiteral struct {
	Token token.Token
	Block *BlockStatement
}

func (n *FunctionLiteral) Node()       {}
func (n *FunctionLiteral) Literal()    {}
func (n *FunctionLiteral) Expression() {}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (a *ArrayLiteral) Node()       {}
func (a *ArrayLiteral) Literal()    {}
func (a *ArrayLiteral) Expression() {}

type ObjectLiteral struct {
	Token    token.Token
	Elements map[Expression]Expression
}

func (o *ObjectLiteral) Node()       {}
func (o *ObjectLiteral) Literal()    {}
func (o *ObjectLiteral) Expression() {}

type TemplateLiteral struct {
	Expressions []Expression
	Components  []token.Token
}

func (t *TemplateLiteral) Node()             {}
func (t *TemplateLiteral) Literal()          {}
func (t *TemplateLiteral) Expression()       {}
func (t *TemplateLiteral) CommandComponent() {}
