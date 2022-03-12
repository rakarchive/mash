package ast

type Program struct {
	Statements []Statement
}

func (p *Program) Node() {}
