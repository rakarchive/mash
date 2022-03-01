package ast

import "github.com/raklaptudirm/mash/pkg/token"

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

func (n *StringLiteral) Node()       {}
func (n *StringLiteral) Literal()    {}
func (n *StringLiteral) Expression() {}

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
