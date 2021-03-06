package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type IR struct {
	Funcs     map[string]*Func
	Nodes     []Node
	Variables []*Variable
}

type FuncParam struct {
	Name  string
	VarID int
	Type  types.Type
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

type empty struct{}

type Builder struct {
	Scope     *ScopeStack
	Variables *[]*Variable
	Funcs     map[string]*Func
	TopLevel  []Node

	fs              FS
	alreadyImported map[string]empty
}

func NewBuilder(fs FS) *Builder {
	vars := make([]*Variable, 0)
	return &Builder{
		Scope:           NewScopeStack(),
		Variables:       &vars,
		TopLevel:        make([]Node, 0),
		fs:              fs,
		alreadyImported: make(map[string]empty),
	}
}

type ScopeStack struct {
	scopes   []*Scope
	scopecnt map[ScopeType]int
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
	}
	return s
}

type ScopeType int

const (
	ScopeTypeGlobal ScopeType = iota
	ScopeTypeFunction
	ScopeTypeIf
	ScopeTypeWhile
	ScopeTypeSwitch
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
	ID   int
	Name string
	Type types.Type
}
