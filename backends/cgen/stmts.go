package cgen

import (
	"github.com/Nv7-Github/gold/ir"
)

func (c *CGen) addNode(s ir.Node) (*Value, error) {
	switch n := s.(type) {
	case *ir.CallNode:
		// Build params
		switch call := n.Call.(type) {
		case *ir.PrintStmt:
			return c.addPrint(call)
		default:
			return nil, s.Pos().Error("unknown call node: %T", c)
		}

	case *ir.Const:
		return c.addConst(n)

	default:
		return nil, s.Pos().Error("unknown node: %T", n)
	}
}
