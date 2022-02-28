package ast

type Node interface {
	Node()
}

type Expression interface {
	Node
	Expression()
}

type Command interface {
	Node
	Command()
}