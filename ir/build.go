package ir

import (
	"errors"

	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type nodeBuilder struct {
	ParamTyps []types.Type
	Build     func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error)
}

var builders = make(map[string]nodeBuilder)

func (b *Builder) Build(p *parser.Parser) (*IR, error) {
	b.alreadyImported[p.Filename()] = empty{} // Already imported

	// Function pass
	err := b.functionPass(p)
	if err != nil {
		return nil, err
	}
	b.Scope.PushScope(NewScope(ScopeTypeGlobal))

	// Build top level
	out := make([]Node, len(p.Nodes))
	for i, stmt := range p.Nodes {
		out[i], err = b.buildNode(stmt)
		if err != nil {
			return nil, err
		}
	}

	if b.Scope.Curr().Type != ScopeTypeGlobal {
		return nil, errors.New("global scope not closed, missing \"end\"?")
	}
	b.Scope.Pop()

	return &IR{
		Funcs:     b.Funcs,
		Nodes:     append(b.TopLevel, out...),
		Variables: *b.Variables,
	}, nil
}

func (b *Builder) buildNode(node parser.Node, inexpr ...bool) (Node, error) {
	switch n := node.(type) {
	case *parser.CallStmt:
		return b.buildCall(n)

	case *parser.Const:
		if n.IsIdentifier && len(inexpr) > 0 { // Variable
			v, exists := b.Scope.GetVar(n.Val.(string))
			if !exists {
				return nil, node.Pos().Error("unknown variable: %s", n.Val.(string))
			}
			return &VariableExpr{
				pos:      n.Pos(),
				typ:      v.Type,
				Variable: v.ID,
			}, nil
		}
		return &Const{
			pos:          n.Pos(),
			typ:          n.Type,
			Value:        n.Val,
			IsIdentifier: n.IsIdentifier,
		}, nil

	case *parser.NotExpr:
		v, err := b.buildNode(n.Val, true)
		if err != nil {
			return nil, err
		}
		if !types.BOOL.Equal(v.Type()) {
			return nil, n.Pos().Error("invalid type for not: %s, expected bool", v.Type())
		}
		return &Not{
			pos: n.Pos(),
			Val: v,
		}, nil

	case *parser.AssignStmt:
		return b.buildAssignStmt(n)

	case *parser.BlockStmt:
		return b.buildBlock(n)

	case *parser.BinaryExpr:
		return b.buildBinaryExpr(n)

	case *parser.UnaryExpr:
		return b.buildUnaryExpr(n)

	case *parser.IndexExpr:
		return b.buildIndexExpr(n)

	default:
		return nil, n.Pos().Error("unknown node type: %T", n)
	}
}

type Const struct {
	pos          *tokenizer.Pos
	typ          types.Type
	Value        interface{}
	IsIdentifier bool
}

func (c *Const) Type() types.Type {
	return c.typ
}

func (c *Const) Pos() *tokenizer.Pos {
	return c.pos
}

type Not struct {
	pos *tokenizer.Pos
	Val Node
}

func (n *Not) Pos() *tokenizer.Pos {
	return n.pos
}

func (n *Not) Type() types.Type {
	return types.BOOL
}

type VariableExpr struct {
	pos      *tokenizer.Pos
	Variable int
	typ      types.Type
}

func (v *VariableExpr) Pos() *tokenizer.Pos { return v.pos }
func (v *VariableExpr) Type() types.Type    { return v.typ }

func (b *Builder) buildCall(n *parser.CallStmt) (Node, error) {
	// Get builder
	bld, exists := builders[n.Fn]
	if !exists {
		return nil, n.Pos().Error("unknown function: %s", n.Fn)
	}

	// Build args
	args := make([]Node, len(n.Args))
	var err error
	for i, arg := range n.Args {
		args[i], err = b.buildNode(arg)
		if err != nil {
			return nil, err
		}
	}

	// Check types
	if err := MatchTypes(n.Pos(), args, bld.ParamTyps); err != nil {
		return nil, err
	}

	c, err := bld.Build(b, n.Pos(), args)
	if err != nil {
		return nil, err
	}
	return &CallNode{
		Call: c,
		pos:  n.Pos(),
	}, nil
}
