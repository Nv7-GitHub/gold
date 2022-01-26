package parser

import (
	"github.com/Nv7-Github/gold/tokenizer"
)

type BlockStmt struct {
	*BasicNode

	Fn    string
	Args  []Node
	Stmts []Node
}

func (p *Parser) parseBlock(fn string, args []Node) (Node, error) {
	stmts := make([]Node, 0)
	ps := p.tok.CurrTok().Pos // BlockStart
	if !p.tok.Eat() {
		return nil, p.getError(ps, "expected block end")
	}

	for p.tok.CurrTok().Type != tokenizer.Operation && (p.tok.CurrTok().Value != tokenizer.BlockEnd) {
		stm, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stm)
	}
	// Eat end
	p.tok.Eat()

	return &BlockStmt{
		BasicNode: &BasicNode{
			pos: ps,
		},
		Fn:    fn,
		Args:  args,
		Stmts: stmts,
	}, nil
}
