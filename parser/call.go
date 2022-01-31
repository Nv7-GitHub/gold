package parser

import (
	"github.com/Nv7-Github/gold/tokenizer"
)

type CallStmt struct {
	*BasicNode

	Fn   string
	Args []Node
}

func (p *Parser) parseCall() (Node, error) {
	fn := p.tok.CurrTok().Value
	ps := p.tok.CurrTok().Pos
	if !p.tok.Eat() {
		return nil, p.getError(ps, "expected \";\"")
	}

	args := make([]Node, 0)
	for p.tok.CurrTok().Type != tokenizer.End {
		if p.tok.CurrTok().Type == tokenizer.Operation {
			switch p.tok.CurrTok().Value {
			case string(tokenizer.Assign):
				return p.parseAssign(&CallStmt{
					BasicNode: &BasicNode{
						pos: ps,
					},
					Fn:   fn,
					Args: args,
				})

			case string(tokenizer.BlockStart):
				return p.parseBlock(fn, args)

			default:
				return nil, p.getError(p.tok.CurrTok().Pos, "unexpected token: %s", p.tok.CurrTok().Value)
			}
		}

		if p.tok.IsEnd() {
			return nil, p.getError(p.tok.CurrPos(), "expected \";\"")
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)
	}
	p.tok.Eat() // End

	return &CallStmt{
		BasicNode: &BasicNode{
			pos: ps,
		},
		Fn:   fn,
		Args: args,
	}, nil
}
