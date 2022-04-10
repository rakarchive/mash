package ast

import "laptudirm.com/x/mash/pkg/token"

type Command interface {
	Node
	Command()
}

type LogicalCommand struct {
	Left     Command
	Operator token.Token
	Right    Command
}

func (l *LogicalCommand) Node()    {}
func (l *LogicalCommand) Command() {}

type UnaryCommand struct {
	Operator token.Token
	Right    Command
}

func (u *UnaryCommand) Node()    {}
func (u *UnaryCommand) Command() {}

type BinaryCommand struct {
	Left     Command
	Operator token.Token
	Right    Command
}

func (b *BinaryCommand) Node()    {}
func (b *BinaryCommand) Command() {}

type LiteralCommand struct {
	Cmd  token.Token
	Args []token.Token
}

func (l *LiteralCommand) Node()    {}
func (l *LiteralCommand) Command() {}
