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
	case val.Type().Equal(types.ARRAY):
		if !ind.Type().Equal(types.INT) {
			return nil, ind.Pos().Error("index must be integer")
		}
		typ = val.Type().(*types.ArrayType).ElemType

	case val.Type().Equal(types.STRING):
		if !ind.Type().Equal(types.INT) {
			return nil, ind.Pos().Error("index must be integer")
		}
		typ = types.STRING

	case val.Type().Equal(types.MAP):
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
