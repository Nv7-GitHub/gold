package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type WhileStmt struct {
	Cond Node
}

func (w *WhileStmt) Type() types.Type { return types.NULL }

func init() {
	builders["while"] = nodeBuilder{
		ParamTyps: []types.Type{types.BOOL},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			b.Scope.PushScope(NewScope(ScopeTypeWhile))
			return &WhileStmt{
				Cond: args[0],
			}, nil
		},
	}
}
