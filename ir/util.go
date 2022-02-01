package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

func MatchTypes(pos *tokenizer.Pos, args []Node, typs []types.Type) error {
	if len(typs) > 0 && typs[len(typs)-1] == types.VARIADIC { // last type is variadic, check up to that
		if len(args) < len(typs)-1 {
			return pos.Error("wrong number of arguments: expected at least %d, got %d", len(typs)-1, len(args))
		}
		for i, t := range typs {
			if t.Equal(types.VARIADIC) {
				return nil
			}

			if !t.Equal(args[i].Type()) {
				return args[i].Pos().Error("wrong argument type: expected %s, got %s", typs[i], args[i].Type())
			}
		}
	}
	if len(args) != len(typs) {
		return pos.Error("wrong number of arguments: expected %d, got %d", len(typs), len(args))
	}
	for i, arg := range args {
		if !arg.Type().Equal(typs[i]) {
			return arg.Pos().Error("wrong argument type: expected %s, got %s", typs[i], arg.Type())
		}
	}
	return nil
}
