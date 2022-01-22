package parser

import (
	"errors"
	"fmt"

	"github.com/Nv7-Github/gold/tokenizer"
)

type Parser struct {
	tok   *tokenizer.Tokenizer
	Nodes []Node
}

func NewParser(tok *tokenizer.Tokenizer) *Parser {
	return &Parser{
		tok: tok,
	}
}

func (p *Parser) Parse() error {
	for !p.tok.IsEnd() {
		n, err := p.parseStmt()
		if err != nil {
			return fmt.Errorf("%s: %s", p.tok.CurrTok().Pos, err)
		}
		p.Nodes = append(p.Nodes, n)
	}
	return nil
}

func (p *Parser) parseExpr() (Expression, error) {
	switch p.tok.CurrTok().Type {
	case tokenizer.NumberLiteral:
		return p.parseNumberLiteral()

	case tokenizer.StringLiteral:
		return p.parseStringLiteral()

	case tokenizer.Identifier:
		return p.parseIdentifier()

	default:
		return nil, errors.New("unknown token")
	}
}

func (p *Parser) parseStmt() (Statement, error) {
	// Get instruction
	tok := p.tok.CurrTok()
	if tok.Type != tokenizer.Identifier {
		return nil, errors.New("expected instruction")
	}
	p.tok.Eat()

	// TODO: Get args, check types, create node
	return nil, errors.New("statements not implemented yet")
}
