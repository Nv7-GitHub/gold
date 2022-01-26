package parser

import (
	"github.com/Nv7-Github/gold/tokenizer"
)

type BinaryExpr struct {
	*BasicNode

	Lhs Node
	Rhs Node
	Op  tokenizer.Op
}

func (p *Parser) parseBinaryExpr() (Node, error) {
	ps := p.tok.CurrTok().Pos
	p.tok.Eat() // get rid of LParen
	lhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Get op
	tok := p.tok.CurrTok()
	if tok.Type != tokenizer.Operator {
		return nil, p.getError(tok.Pos, "expected operator")
	}
	if !p.tok.Eat() {
		return nil, p.getError(p.tok.CurrPos(), "expected rhs")
	}

	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Eat last parenthesis
	if !p.tok.Eat() {
		return nil, p.getError(tok.Pos, "expected \")\"")
	}

	return &BinaryExpr{
		BasicNode: &BasicNode{
			pos: ps,
		},
		Lhs: lhs,
		Rhs: rhs,
		Op:  tokenizer.Op([]rune(tok.Value)[0]),
	}, nil
}
