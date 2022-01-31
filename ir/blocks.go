package ir

import "github.com/Nv7-Github/gold/parser"

func (b *Builder) buildBlock(n *parser.BlockStmt) (Node, error) {
	// Get builder
	bld, exists := builders[n.Fn]
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
	c, err := bld.Build(b, n.Pos(), args)
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

	// Pop scope
	b.Scope.Pop()

	return &BlockNode{
		Call: c,
		Body: body,
		pos:  n.Pos(),
	}, nil
}
