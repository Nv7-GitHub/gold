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
	p.tok.Eat()

	args := make([]Node, 0)
	for p.tok.CurrTok().Type != tokenizer.End && (p.tok.CurrTok().Type != tokenizer.Operation) {
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
