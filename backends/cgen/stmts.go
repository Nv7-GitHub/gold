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

		case *ir.CallStmt:
			return c.addCall(call)

		case *ir.ReturnStmt:
			return c.addReturn(call)

		case *ir.DefCall:
			return c.addDef(call)

		default:
			return nil, s.Pos().Error("unknown call node: %T", call)
		}

	case *ir.Const:
		return c.addConst(n)

	case *ir.VariableExpr:
		return c.addVarExpr(n)

	case *ir.MathExpr:
		return c.addMathExpr(n)

	case *ir.AssignStmt:
		return c.addAssign(n)

	default:
		return nil, s.Pos().Error("unknown node: %T", n)
	}
}
