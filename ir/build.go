package ir

import (
	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type nodeBuilder struct {
	ParamTyps []types.Type
	Build     func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error)
}

var builders = make(map[string]nodeBuilder)

func (b *Builder) Build(p *parser.Parser) ([]Node, error) {
	out := make([]Node, len(p.Nodes))
	var err error
	for i, stmt := range p.Nodes {
		out[i], err = b.buildNode(stmt)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
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
				pos: n.Pos(),
				typ: v.Type,
			}, nil
		}
		return &Const{
			pos:          n.Pos(),
			typ:          n.Type,
			Value:        n.Val,
			IsIdentifier: n.IsIdentifier,
		}, nil

	case *parser.AssignStmt:
		return b.buildAssignStmt(n)

	case *parser.BlockStmt:
		return b.buildBlock(n)

	case *parser.BinaryExpr:
		return b.buildBinaryExpr(n)

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

type VariableExpr struct {
	pos  *tokenizer.Pos
	Name string
	typ  types.Type
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
