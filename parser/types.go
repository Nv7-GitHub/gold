package parser

import (
	"github.com/Nv7-Github/gold/types"

	"github.com/Nv7-Github/gold/tokenizer"
)

type Node interface {
	Pos() *tokenizer.Pos
}

type BasicNode struct {
	pos *tokenizer.Pos
}

func (b *BasicNode) Pos() *tokenizer.Pos {
	return b.pos
}

type Statement interface {
	Node
}

type Expression interface {
	Node

	Type() types.Type
}
