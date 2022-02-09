package cgen

import (
	"fmt"

	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/types"
)

func (c *CGen) addWhile(s *ir.WhileStmt) (*Value, error) {
	cond, err := c.addNode(s.Cond)
	if err != nil {
		return nil, err
	}

	body, err := c.BuildStmts(s.Body)
	if err != nil {
		return nil, err
	}

	return &Value{
		Setup:    cond.Setup,
		Destruct: cond.Destruct,
		Code:     fmt.Sprintf("while (%s) {\n%s\n}", cond.Code, Indent(body)),
		Type:     types.NULL,
	}, nil
}

func (c *CGen) addIf(s *ir.IfStmt) (*Value, error) {
	cond, err := c.addNode(s.Cond)
	if err != nil {
		return nil, err
	}

	body, err := c.BuildStmts(s.Body)
	if err != nil {
		return nil, err
	}
	elseCode := ""
	if s.Else != nil {
		els, err := c.BuildStmts(s.Else)
		if err != nil {
			return nil, err
		}
		elseCode = fmt.Sprintf("else {\n%s\n}", Indent(els))
	}

	return &Value{
		Setup:    cond.Setup,
		Destruct: cond.Destruct,
		Code:     fmt.Sprintf("if (%s) {\n%s\n}%s", cond.Code, Indent(body), elseCode),
		Type:     types.NULL,
	}, nil
}
