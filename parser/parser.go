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

type ExprStmt struct {
	Expr Node
}

func (e *ExprStmt) Pos() *tokenizer.Pos {
	return e.Expr.Pos()
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

func (p *Parser) parseExpr() (Node, error) {
	switch p.tok.CurrTok().Type {
	case tokenizer.NumberLiteral:
		return p.parseNumberLiteral()

	case tokenizer.StringLiteral:
		return p.parseStringLiteral()

	case tokenizer.BoolLiteral:
		return p.parseBoolLiteral()

	case tokenizer.Identifier:
		return p.parseIdentifier()

	case tokenizer.Parenthesis:
		if p.tok.CurrTok().Value == string(tokenizer.LParen) {
			return p.parseOp()
		}
		return nil, p.getError(p.tok.CurrTok().Pos, "unknown token: %s", p.tok.CurrTok().Value)

	default:
		return nil, p.getError(p.tok.CurrTok().Pos, "unknown token: %s", p.tok.CurrTok().Value)
	}
}

func (p *Parser) parseStmt() (Node, error) {
	// Get instruction
	tok := p.tok.CurrTok()
	switch tok.Type {
	case tokenizer.Identifier:
		return p.parseCall()

	case tokenizer.Parenthesis:
		// (a + b) => c;
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		if p.tok.CurrTok().Type == tokenizer.Operation && p.tok.CurrTok().Value == tokenizer.Assign {
			return p.parseAssign(expr)
		} else {
			return &ExprStmt{
				Expr: expr,
			}, nil
		}

	default:
		return nil, p.getError(tok.Pos, "unexpected token: %s", tok.Value)
	}
}
