package parser

import "github.com/Nv7-Github/gold/tokenizer"

type AssignStmt struct {
	*BasicNode

	Value    Node
	Variable string
}

func (p *Parser) parseAssign(expr Node) (Node, error) {
	ps := p.tok.CurrTok().Pos
	p.tok.Eat()
	if p.tok.CurrTok().Type != tokenizer.Identifier {
		return nil, p.getError(p.tok.CurrTok().Pos, "expected identifier")
	}
	varname := p.tok.CurrTok().Value
	p.tok.Eat()

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
		Variable: varname,
		Value:    expr,
	}, nil
}
