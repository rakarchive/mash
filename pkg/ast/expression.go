package ast

type Expression interface {
	Node
	Expression()
}
