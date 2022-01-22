package parser

import (
	"strconv"

	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type Const struct {
	*BasicNode

	typ          types.Type
	Val          interface{}
	IsIdentifier bool
}

func (c *Const) Type() types.Type {
	return c.typ
}

func (p *Parser) parseNumberLiteral() (Expression, error) {
	tok := p.tok.CurrTok()
	if tok.Type != tokenizer.NumberLiteral {
		return nil, p.getError(tok.Pos, "expected number literal")
	}
	p.tok.Eat()
	v, err := strconv.Atoi(tok.Value)
	if err != nil {
		return nil, err
	}
	return &Const{
		BasicNode: &BasicNode{
			pos: p.tok.CurrTok().Pos,
		},
		typ: types.INT,
		Val: v,
	}, nil
}

func (p *Parser) parseStringLiteral() (Expression, error) {
	tok := p.tok.CurrTok()
	if tok.Type != tokenizer.StringLiteral {
		return nil, p.getError(tok.Pos, "expected string literal")
	}
	p.tok.Eat()
	return &Const{
		BasicNode: &BasicNode{
			pos: p.tok.CurrTok().Pos,
		},
		typ: types.STRING,
		Val: tok.Value,
	}, nil
}

func (p *Parser) parseIdentifier() (Expression, error) {
	tok := p.tok.CurrTok()
	if tok.Type != tokenizer.Identifier {
		return nil, p.getError(tok.Pos, "expected identifier")
	}
	p.tok.Eat()
	return &Const{
		BasicNode: &BasicNode{
			pos: p.tok.CurrTok().Pos,
		},
		typ:          types.STRING,
		Val:          tok.Value,
		IsIdentifier: true,
	}, nil
}
