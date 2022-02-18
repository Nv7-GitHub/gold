package ir

import (
	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type blockBuilder struct {
	ParamTyps []types.Type
	Init      func(b *Builder, pos *tokenizer.Pos, args []Node) (Block, error)
	Build     func(b *Builder, pos *tokenizer.Pos, blk Block, stmts []Node) error
}

var blockBuilders = make(map[string]blockBuilder)

func (b *Builder) buildBlock(n *parser.BlockStmt) (Node, error) {
	// Get builder
	bld, exists := blockBuilders[n.Fn]
	if !exists {
		return nil, n.Pos().Error("unknown block type: %s", n.Fn)
	}

	// Call
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

	// Build
	c, err := bld.Init(b, n.Pos(), args)
	if err != nil {
		return nil, err
	}

	// Body
	body := make([]Node, len(n.Stmts))
	for i, stmt := range n.Stmts {
		body[i], err = b.buildNode(stmt)
		if err != nil {
			return nil, err
		}
	}

	// Finalize
	err = bld.Build(b, n.Pos(), c, body)
	if err != nil {
		return nil, err
	}

	return &BlockNode{
		Block: c,
		pos:   n.Pos(),
	}, nil
}
