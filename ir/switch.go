package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type SwitchStmt struct {
	Cond  Node
	Cases []*SwitchCase
}

type SwitchCase struct {
	Cond *Const
	Body []Node
}

func init() {
	blockBuilders["case"] = blockBuilder{
		ParamTyps: []types.Type{types.ANY},
		Init: func(b *Builder, pos *tokenizer.Pos, args []Node) (Block, error) {
			cond, ok := args[0].(*Const)
			if !ok {
				return nil, args[0].Pos().Error("case condition must be a constant")
			}
			return &SwitchCase{
				Cond: cond,
			}, nil
		},
		Build: func(b *Builder, pos *tokenizer.Pos, blk Block, stmts []Node) error {
			blk.(*SwitchCase).Body = stmts
			return nil
		},
	}

	blockBuilders["switch"] = blockBuilder{
		ParamTyps: []types.Type{types.ANY},
		Init: func(b *Builder, pos *tokenizer.Pos, args []Node) (Block, error) {
			b.Scope.PushScope(NewScope(ScopeTypeSwitch))
			return &SwitchStmt{
				Cond: args[0],
			}, nil
		},
		Build: func(b *Builder, pos *tokenizer.Pos, blk Block, stmts []Node) error {
			sw := blk.(*SwitchStmt)
			cases := make([]*SwitchCase, len(stmts))
			for i, stmt := range stmts {
				// Convert to case
				blk, ok := stmt.(*BlockNode)
				if !ok {
					return stmt.Pos().Error("switch case must only have case statements within")
				}
				cs, ok := blk.Block.(*SwitchCase)
				if !ok {
					return blk.Pos().Error("switch case must only have case statements within")
				}

				// Check type
				if !cs.Cond.Type().Equal(sw.Cond.Type()) {
					return blk.Pos().Error("case condition must have the same type as switch condition")
				}

				// Save
				cases[i] = cs
			}
			sw.Cases = cases
			b.Scope.Pop()
			return nil
		},
	}
}
