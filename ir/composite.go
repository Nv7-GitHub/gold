package ir

import (
	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type IndexExpr struct {
	pos *tokenizer.Pos

	Value Node
	Index Node

	typ types.Type
}

func (i *IndexExpr) Type() types.Type    { return i.typ }
func (i *IndexExpr) Pos() *tokenizer.Pos { return i.pos }

func (b *Builder) buildIndexExpr(s *parser.IndexExpr) (Node, error) {
	val, err := b.buildNode(s.Value, true)
	if err != nil {
		return nil, err
	}
	ind, err := b.buildNode(s.Index, true)
	if err != nil {
		return nil, err
	}

	var typ types.Type
	switch {
	case types.ARRAY.Equal(val.Type()):
		if !ind.Type().Equal(types.INT) {
			return nil, ind.Pos().Error("index must be integer")
		}
		typ = val.Type().(*types.ArrayType).ElemType

	case types.STRING.Equal(val.Type()):
		if !ind.Type().Equal(types.INT) {
			return nil, ind.Pos().Error("index must be integer")
		}
		typ = types.STRING

	case types.MAP.Equal(val.Type()):
		if !ind.Type().Equal(val.Type().(*types.MapType).KeyType) {
			return nil, ind.Pos().Error("expected index type %s, got index type %s", val.Type().(*types.MapType).KeyType.String(), ind.Type().String())
		}
		typ = val.Type().(*types.MapType).ValType

	default:
		return nil, val.Pos().Error("indexing not supported for type %s", val.Type().String())
	}

	return &IndexExpr{
		pos:   s.Pos(),
		Value: val,
		Index: ind,
		typ:   typ,
	}, nil
}

type AppendStmt struct {
	Array Node
	Val   Node
}

func (a *AppendStmt) Type() types.Type { return types.NULL }

type GrowStmt struct {
	Array Node
	Size  Node
}

func (g *GrowStmt) Type() types.Type { return types.NULL }

type SetStmt struct {
	Map   Node
	Key   Node
	Value Node
}

func (s *SetStmt) Type() types.Type { return types.NULL }

type GetStmt struct {
	Map Node
	Key Node

	typ types.Type
}

func (g *GetStmt) Type() types.Type { return g.typ }

func init() {
	builders["append"] = nodeBuilder{
		ParamTyps: []types.Type{types.ARRAY, types.ANY},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			elType := args[0].Type().(*types.ArrayType).ElemType
			if !args[1].Type().Equal(elType) {
				return nil, args[1].Pos().Error("expected type %s, got type %s", elType.String(), args[1].Type().String())
			}

			return &AppendStmt{
				Array: args[0],
				Val:   args[1],
			}, nil
		},
	}

	builders["grow"] = nodeBuilder{
		ParamTyps: []types.Type{types.ARRAY, types.INT},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			return &GrowStmt{
				Array: args[0],
				Size:  args[1],
			}, nil
		},
	}

	builders["set"] = nodeBuilder{
		ParamTyps: []types.Type{types.MAP, types.ANY, types.ANY},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			mapTyp := args[0].Type().(*types.MapType)
			if !args[1].Type().Equal(mapTyp.KeyType) {
				return nil, args[1].Pos().Error("expected type %s, got type %s", mapTyp.KeyType.String(), args[1].Type().String())
			}
			if !args[2].Type().Equal(mapTyp.ValType) {
				return nil, args[2].Pos().Error("expected type %s, got type %s", mapTyp.ValType.String(), args[2].Type().String())
			}
			return &SetStmt{
				Map:   args[0],
				Key:   args[1],
				Value: args[2],
			}, nil
		},
	}

	builders["get"] = nodeBuilder{
		ParamTyps: []types.Type{types.MAP, types.ANY},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			mapTyp := args[0].Type().(*types.MapType)
			if !args[1].Type().Equal(mapTyp.KeyType) {
				return nil, args[1].Pos().Error("expected type %s, got type %s", mapTyp.KeyType.String(), args[1].Type().String())
			}
			return &GetStmt{
				Map: args[0],
				Key: args[1],
				typ: mapTyp.ValType,
			}, nil
		},
	}
}
