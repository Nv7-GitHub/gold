package ir

import (
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

func MatchTypes(pos *tokenizer.Pos, args []Node, typs []types.Type) error {
	if len(args) != len(typs) {
		return pos.Error("wrong number of arguments: expected %d, got %d", len(typs), len(args))
	}
	for i, arg := range args {
		if !arg.Type().Equal(typs[i]) {
			return pos.Error("wrong argument type: expected %s, got %s", typs[i], arg.Type())
		}
	}
	return nil
}
