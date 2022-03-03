package ast

import "github.com/raklaptudirm/mash/pkg/token"

type Expression interface {
	Node
	Expression()
}

type Assignable interface {
	Expression
	Assignable()
}

type AssignExpression struct {
	Left     Assignable
	Operator token.Token
	Right    Expression
}

func (a *AssignExpression) Node()       {}
func (a *AssignExpression) Expression() {}

type LogicalExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (l *LogicalExpression) Node()       {}
func (l *LogicalExpression) Expression() {}

type BinaryExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (b *BinaryExpression) Node()       {}
func (b *BinaryExpression) Expression() {}

type UnaryExpression struct {
	Operator token.Token
	Right    Expression
}

func (u *UnaryExpression) Node()       {}
func (u *UnaryExpression) Expression() {}

type GroupExpression struct {
	Right Expression
}

func (g *GroupExpression) Node()       {}
func (g *GroupExpression) Expression() {}

type CallExpression struct {
	Callee      Expression
	Parenthesis token.Token
	Arguments   []Expression
}

func (c *CallExpression) Node()       {}
func (c *CallExpression) Expression() {}

type GetExpression struct {
	Name Expression
	Expr Expression
}

func (g *GetExpression) Node()       {}
func (g *GetExpression) Expression() {}
func (g *GetExpression) Assignable() {}

type VariableExpression struct {
	Name token.Token
}

func (v *VariableExpression) Node()       {}
func (v *VariableExpression) Expression() {}
func (v *VariableExpression) Assignable() {}

type LiteralExpression struct {
	Literal Literal
}

func (l *LiteralExpression) Node()       {}
func (l *LiteralExpression) Expression() {}
