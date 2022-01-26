package parser

import "github.com/Nv7-Github/gold/tokenizer"

type BinaryExpr struct {
	*BasicNode

	Lhs Node
	Rhs Node
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

	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Eat last parenthesis
	p.tok.Eat()

	return &BinaryExpr{
		BasicNode: &BasicNode{
			pos: ps,
		},
		Lhs: lhs,
		Rhs: rhs,
	}, nil
}
