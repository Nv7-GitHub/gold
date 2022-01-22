package parser

import (
	"github.com/Nv7-Github/gold/tokenizer"
)

func (p *Parser) getStmt(pos *tokenizer.Pos, name string, args []Expression) (Statement, error) {
	// TODO: Get instr parser, check types, parse instr
	return nil, p.getError(pos, "nodes not implemented")
}
