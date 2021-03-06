package cgen

import (
	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/types"
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
			return c.addDef(n.Pos(), call)

		case *ir.StringCast:
			return c.addStringCast(call)

		case *ir.AppendStmt:
			return c.addAppendStmt(call)

		case *ir.GrowStmt:
			return c.addGrowStmt(call)

		case *ir.LengthStmt:
			return c.addLen(call)

		case *ir.StringConcat:
			return c.addConcat(call)

		case *ir.ImportCall:
			return &Value{
				Type: types.NULL,
			}, nil

		case *ir.ExitStmt:
			return c.addExit(call)

		default:
			return nil, s.Pos().Error("unknown call node: %T", call)
		}

	case *ir.BlockNode:
		switch blk := n.Block.(type) {
		case *ir.WhileStmt:
			return c.addWhile(blk)

		case *ir.IfStmt:
			return c.addIf(blk)

		case *ir.SwitchStmt:
			return c.addSwitch(blk)

		default:
			return nil, s.Pos().Error("unknown block node: %T", blk)
		}

	case *ir.Const:
		return c.addConst(n)

	case *ir.VariableExpr:
		return c.addVarExpr(n)

	case *ir.MathExpr:
		return c.addMathExpr(n)

	case *ir.AssignStmt:
		return c.addAssign(n)

	case *ir.ComparisonExpr:
		return c.addComparison(n)

	case *ir.IndexExpr:
		return c.addIndexExpr(n)

	case *ir.StringEq:
		return c.addStringEq(n)

	case *ir.Not:
		v, err := c.addNode(n.Val)
		if err != nil {
			return nil, err
		}
		return &Value{
			Setup:    v.Setup,
			Destruct: v.Destruct,
			Type:     types.BOOL,
			Code:     "!" + v.Code,
		}, nil

	default:
		return nil, s.Pos().Error("unknown node: %T", n)
	}
}
