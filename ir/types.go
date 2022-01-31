package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type Node interface {
	Pos() *tokenizer.Pos
	Type() types.Type
}

type Call interface {
	Type() types.Type
}

type CallNode struct {
	Call Call
	pos  *tokenizer.Pos
}

func (c *CallNode) Pos() *tokenizer.Pos {
	return c.pos
}

func (c *CallNode) Type() types.Type {
	return c.Call.Type()
}

type BlockNode struct {
	Call Call
	Body []Node

	pos *tokenizer.Pos
}

func (b *BlockNode) Pos() *tokenizer.Pos {
	return b.pos
}

func (b *BlockNode) Type() types.Type {
	return types.NULL
}

type Builder struct {
	Scope *ScopeStack
}

func NewBuilder() *Builder {
	return &Builder{
		Scope: NewScopeStack(),
	}
}

type ScopeStack struct {
	scopes   []*Scope
	scopecnt map[ScopeType]int
	vars     map[string]*Variable
}

func NewScopeStack() *ScopeStack {
	s := &ScopeStack{
		scopes:   make([]*Scope, 0, 1),
		scopecnt: make(map[ScopeType]int, 1),
		vars:     make(map[string]*Variable),
	}
	s.PushScope(NewScope(ScopeTypeGlobal))
	return s
}

type ScopeType int

const (
	ScopeTypeGlobal ScopeType = iota
	ScopeTypeFunction
	ScopeTypeIf
	ScopeTypeWhile
)

type Scope struct {
	Type      ScopeType
	Variables map[string]*Variable

	parent *ScopeStack
}

func NewScope(typ ScopeType) *Scope {
	return &Scope{
		Type:      typ,
		Variables: make(map[string]*Variable),
	}
}

type Variable struct {
	Name string
	Type types.Type
}
