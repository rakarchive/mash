package ast

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
	ElseBlock Statement
}

func (i *IfStatement) Node()      {}
func (i *IfStatement) Statement() {}

type ForStatement struct {
	Condition Expression
	BlockStmt *BlockStatement
}

func (f *ForStatement) Node()      {}
func (f *ForStatement) Statement() {}

type LetStatement struct {
	Expression Expression
}

func (l *LetStatement) Node()      {}
func (l *LetStatement) Statement() {}

type CmdStatement struct {
	Command Command
}

func (c *CmdStatement) Node()      {}
func (c *CmdStatement) Statement() {}
