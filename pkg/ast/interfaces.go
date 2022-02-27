package ast

type Node interface {
	Node()
}

type Statement interface {
	Node
	Statement()
}

type Expression interface {
	Node
	Expression()
}

type Command interface {
	Node
	Command()
}
