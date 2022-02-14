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

type UnaryExpr struct {
	*BasicNode

	Val Node
}

func (p *Parser) parseOp() (Node, error) {
	ps := p.tok.CurrTok().Pos
	p.tok.Eat() // get rid of LParen
	lhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Get op
	tok := p.tok.CurrTok()
	if tok.Type != tokenizer.Operator {
		if tok.Type == tokenizer.Parenthesis && tok.Value == string(tokenizer.RParen) {
			p.tok.Eat()
			return &UnaryExpr{
				BasicNode: &BasicNode{
					pos: ps,
				},
				Val: lhs,
			}, nil
		}
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

type IndexExpr struct {
	*BasicNode

	Index Node
	Value Node
}

func (p *Parser) parseIndex() (Node, error) {
	// Parse [x]
	ps := p.tok.CurrTok().Pos
	p.tok.Eat() // get rid of LBrack
	ind, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	tok := p.tok.CurrTok()
	if tok.Type != tokenizer.Parenthesis && tok.Value != string(tokenizer.RBrack) {
		return nil, p.getError(tok.Pos, "expected \"]\"")
	} else {
		p.tok.Eat()
	}

	// Get value
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &IndexExpr{
		BasicNode: &BasicNode{
			pos: ps,
		},
		Index: ind,
		Value: val,
	}, nil
}
