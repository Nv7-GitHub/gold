package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type PrintStmt struct {
	Arg Node
}

func (p *PrintStmt) Type() types.Type { return types.NULL }

func init() {
	builders["print"] = nodeBuilder{
		ParamTyps: []types.Type{types.STRING},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			return &PrintStmt{
				Arg: args[0],
			}, nil
		},
	}
}
