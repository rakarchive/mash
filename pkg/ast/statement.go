package ast

import "github.com/raklaptudirm/mash/pkg/token"

type Statement interface {
	Node
	Statement()
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) Node()      {}
func (b *BlockStatement) Statement() {}

type IfStatement struct {
	Condition Expression
	BlockStmt *BlockStatement
	ElifBlock []ElifBlock
	ElseBlock *BlockStatement
}

func (i *IfStatement) Node()      {}
func (i *IfStatement) Statement() {}

type ElifBlock struct {
	Condition Expression
	BlockStmt *BlockStatement
}

type ForStatement struct {
	Condition Expression
	BlockStmt *BlockStatement
}

func (f *ForStatement) Node()      {}
func (f *ForStatement) Statement() {}

type LetStatement struct {
	AssignOp   token.TokenType
	Assignable Expression
	Expression Expression
}

func (l *LetStatement) Node()      {}
func (l *LetStatement) Statement() {}

type CmdStatement struct {
	Command Command
}

func (c *CmdStatement) Node()      {}
func (c *CmdStatement) Statement() {}
