package ir

import (
	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type Cast struct {
	Value  Node
	NewTyp types.Type
}

func (c *Cast) Pos() *tokenizer.Pos { return c.Value.Pos() }
func (c *Cast) Type() types.Type    { return c.NewTyp }

type MathExpr struct {
	pos *tokenizer.Pos
	Op  tokenizer.Op
	Lhs Node
	Rhs Node

	typ types.Type
}

func (m *MathExpr) Pos() *tokenizer.Pos { return m.pos }
func (m *MathExpr) Type() types.Type    { return m.typ }

type ComparisonExpr struct {
	pos *tokenizer.Pos
	Op  tokenizer.Op
	Lhs Node
	Rhs Node
}

func (c *ComparisonExpr) Pos() *tokenizer.Pos { return c.pos }
func (c *ComparisonExpr) Type() types.Type    { return types.BOOL }

type LogicalExpr struct {
	pos *tokenizer.Pos
	Op  tokenizer.Op
	Lhs Node
	Rhs Node
}

func (l *LogicalExpr) Pos() *tokenizer.Pos { return l.pos }
func (l *LogicalExpr) Type() types.Type    { return types.BOOL }

func (b *Builder) buildBinaryExpr(n *parser.BinaryExpr) (Node, error) {
	lhs, err := b.buildNode(n.Lhs, true)
	if err != nil {
		return nil, err
	}
	rhs, err := b.buildNode(n.Rhs, true)
	if err != nil {
		return nil, err
	}

	switch n.Op {
	case tokenizer.Add, tokenizer.Sub, tokenizer.Mul, tokenizer.Div, tokenizer.Mod:
		return b.buildMathOp(n, lhs, rhs)

	case tokenizer.Eq, tokenizer.Ne, tokenizer.Gt, tokenizer.Lt:
		return b.buildComparisonOp(n, lhs, rhs)

	case tokenizer.And, tokenizer.Or:
		return b.buildLogicalOp(n, lhs, rhs)

	default:
		return nil, n.Pos().Error("unsupported operator: %s", n.Op)
	}
}

func getCommonNumType(lhs, rhs Node) (Node, Node, types.Type, error) {
	if !lhs.Type().Equal(types.FLOAT) && !lhs.Type().Equal(types.INT) {
		return nil, nil, types.NULL, lhs.Pos().Error("cannot operate on non-numeric types")
	}
	if !rhs.Type().Equal(types.FLOAT) && !rhs.Type().Equal(types.INT) {
		return nil, nil, types.NULL, rhs.Pos().Error("cannot operate on non-numeric types")
	}

	var typ types.Type
	if !lhs.Type().Equal(rhs.Type()) { // Cast to float if not equal, since can't cast to int
		typ = types.FLOAT
		if !lhs.Type().Equal(types.FLOAT) {
			lhs = &Cast{
				Value:  lhs,
				NewTyp: types.FLOAT,
			}
		}

		if !rhs.Type().Equal(types.FLOAT) {
			rhs = &Cast{
				Value:  rhs,
				NewTyp: types.FLOAT,
			}
		}
	} else {
		typ = lhs.Type()
	}

	return lhs, rhs, typ, nil
}

func (b *Builder) buildMathOp(n *parser.BinaryExpr, lhs, rhs Node) (Node, error) {
	var typ types.Type
	var err error
	lhs, rhs, typ, err = getCommonNumType(lhs, rhs)
	if err != nil {
		return nil, err
	}

	return &MathExpr{
		pos: n.Pos(),
		Op:  n.Op,
		Lhs: lhs,
		Rhs: rhs,

		typ: typ,
	}, nil
}

func (b *Builder) buildComparisonOp(n *parser.BinaryExpr, lhs, rhs Node) (Node, error) {
	var err error
	lhs, rhs, _, err = getCommonNumType(lhs, rhs)
	if err != nil {
		return nil, err
	}

	return &ComparisonExpr{
		pos: n.Pos(),
		Op:  n.Op,
		Lhs: lhs,
		Rhs: rhs,
	}, nil
}

func (b *Builder) buildLogicalOp(n *parser.BinaryExpr, lhs, rhs Node) (Node, error) {
	if !lhs.Type().Equal(types.BOOL) {
		return nil, lhs.Pos().Error("cannot perform logical operations on non-bool types")
	}
	if !rhs.Type().Equal(types.BOOL) {
		return nil, rhs.Pos().Error("cannot perform logical operations on non-bool types")
	}

	return &LogicalExpr{
		pos: n.Pos(),
		Op:  n.Op,
		Lhs: lhs,
		Rhs: rhs,
	}, nil
}
