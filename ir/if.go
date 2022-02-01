package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type IfStmt struct {
	Cond Node
	Body []Node
	Else []Node
}

type ElseStmt struct{}

func (e *ElseStmt) Type() types.Type { return types.NULL }

func init() {
	builders["else"] = nodeBuilder{
		ParamTyps: []types.Type{},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			if b.Scope.Curr().Type != ScopeTypeIf {
				return nil, pos.Error("else statement outside of if")
			}
			if b.Scope.Curr().ElsePos != nil {
				return nil, pos.Error("duplicate else statement, original at %s", b.Scope.Curr().ElsePos.String())
			}
			b.Scope.Curr().ElsePos = pos
			return &ElseStmt{}, nil
		},
	}

	blockBuilders["if"] = blockBuilder{
		ParamTyps: []types.Type{types.BOOL},
		Init: func(b *Builder, pos *tokenizer.Pos, args []Node) (Block, error) {
			b.Scope.PushScope(NewScope(ScopeTypeIf))
			return &IfStmt{
				Cond: args[0],
			}, nil
		},
		Build: func(b *Builder, pos *tokenizer.Pos, blk Block, stmts []Node) error {
			var _else []Node
			for i, stmt := range stmts {
				_, ok := stmt.(*CallNode).Call.(*ElseStmt)
				if ok {
					_else = stmts[i+1:]
					stmts = stmts[:i]
					break
				}
			}
			blk.(*IfStmt).Body = stmts
			blk.(*IfStmt).Else = _else
			return nil
		},
	}
}
