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
		ParamTyps: []types.Type{types.ANY},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			val := args[0]
			if !args[0].Type().Equal(types.STRING) {
				fmtc, err := getStringCast(b, pos, args)
				if err != nil {
					return nil, err
				}
				val = &CallNode{
					Call: fmtc,
					pos:  pos,
				}
			}
			return &PrintStmt{
				Arg: val,
			}, nil
		},
	}
}
