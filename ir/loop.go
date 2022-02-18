package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type WhileStmt struct {
	Cond Node
	Body []Node
}

func init() {
	blockBuilders["while"] = blockBuilder{
		ParamTyps: []types.Type{types.BOOL},
		Init: func(b *Builder, pos *tokenizer.Pos, args []Node) (Block, error) {
			b.Scope.PushScope(NewScope(ScopeTypeWhile))
			return &WhileStmt{
				Cond: args[0],
			}, nil
		},
		Build: func(b *Builder, pos *tokenizer.Pos, blk Block, stmts []Node) error {
			blk.(*WhileStmt).Body = stmts
			b.Scope.Pop()
			return nil
		},
	}
}
