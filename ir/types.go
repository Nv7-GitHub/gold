package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type IR struct {
	Funcs map[string]*Func
	Nodes []Node
}

type FuncParam struct {
	Name string
	Type types.Type
}

type Func struct {
	pos *tokenizer.Pos

	Name    string
	Params  []FuncParam
	RetType types.Type
	Body    []Node
}

func (f *Func) Pos() *tokenizer.Pos { return f.pos }

type Node interface {
	Pos() *tokenizer.Pos
	Type() types.Type
}

type Call interface {
	Type() types.Type
}

type Block interface{}

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
	Block Block

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
	Funcs map[string]*Func
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

func (s *ScopeStack) HasScope(t ScopeType) bool {
	cnt, exists := s.scopecnt[t]
	if !exists {
		return false
	}
	return cnt > 0
}

func (s *ScopeStack) GetScopeByType(t ScopeType) *Scope {
	for _, scope := range s.scopes {
		if scope.Type == t {
			return scope
		}
	}
	return nil
}

func NewScopeStack() *ScopeStack {
	s := &ScopeStack{
		scopes:   make([]*Scope, 0, 1),
		scopecnt: make(map[ScopeType]int, 1),
		vars:     make(map[string]*Variable),
	}
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

	// Scope-specific things
	ElsePos  *tokenizer.Pos
	FuncName string

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
