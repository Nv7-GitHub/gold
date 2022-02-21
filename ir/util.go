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
		if !typs[i].Equal(arg.Type()) {
			return arg.Pos().Error("wrong argument type: expected %s, got %s", typs[i], arg.Type())
		}
	}
	return nil
}

// Utility functions in Gold
type StringCast struct {
	Arg Node
}

func (s *StringCast) Type() types.Type { return types.STRING }

func getStringCast(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
	typ := args[0].Type()
	if !typ.Equal(types.FLOAT) && !typ.Equal(types.INT) {
		return nil, args[0].Pos().Error("can only format float and int")
	}
	return &StringCast{
		Arg: args[0],
	}, nil
}

type StringConcat struct {
	Lhs Node
	Rhs Node
}

func (s *StringConcat) Type() types.Type { return types.STRING }

func init() {
	builders["str"] = nodeBuilder{
		ParamTyps: []types.Type{types.ANY},
		Build:     getStringCast,
	}

	builders["concat"] = nodeBuilder{
		ParamTyps: []types.Type{types.STRING, types.STRING},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			return &StringConcat{
				Lhs: args[0],
				Rhs: args[1],
			}, nil
		},
	}
}
