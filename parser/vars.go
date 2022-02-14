package parser

import (
	"github.com/Nv7-Github/gold/tokenizer"
)

type AssignStmt struct {
	*BasicNode

	Value    Node
	Variable Node
}

func (p *Parser) parseAssign(expr Node) (Node, error) {
	ps := p.tok.CurrTok().Pos
	p.tok.Eat()
	var varv Node
	switch p.tok.CurrTok().Type {
	case tokenizer.Identifier:
		var err error
		varv, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}

	case tokenizer.Parenthesis:
		if p.tok.CurrTok().Value == string(tokenizer.LParen) {
			var err error
			varv, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
			break
		}
		fallthrough

	default:
		return nil, p.getError(ps, "expected identifier or parenthesis")
	}

	// Semicolon
	if p.tok.CurrTok().Type != tokenizer.End {
		return nil, p.getError(p.tok.CurrTok().Pos, "expected end")
	} else {
		p.tok.Eat()
	}

	return &AssignStmt{
		BasicNode: &BasicNode{
			pos: ps,
		},
		Variable: varv,
		Value:    expr,
	}, nil
}
