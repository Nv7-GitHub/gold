package parser

import (
	"fmt"

	"github.com/Nv7-Github/gold/tokenizer"
)

type Parser struct {
	tok   *tokenizer.Tokenizer
	Nodes []Node
}

func (p *Parser) getError(pos *tokenizer.Pos, format string, args ...interface{}) error {
	return fmt.Errorf("%s: %s", pos, fmt.Sprintf(format, args...))
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
			return err
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
		return nil, p.getError(p.tok.CurrTok().Pos, "unknown token")
	}
}

func (p *Parser) parseStmt() (Statement, error) {
	// Get instruction
	tok := p.tok.CurrTok()
	if tok.Type != tokenizer.Identifier {
		return nil, p.getError(p.tok.CurrTok().Pos, "expected instruction")
	}
	p.tok.Eat()

	args := []Expression{}
	for !p.tok.IsEnd() {
		if p.tok.CurrTok().Type == tokenizer.End {
			p.tok.Eat()
			break
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)
	}

	return p.getStmt(tok.Pos, tok.Value, args)
}
