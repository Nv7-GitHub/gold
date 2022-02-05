package cgen

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/types"
)

func (c *CGen) addCall(s *ir.CallStmt) (*Value, error) {
	setup := &strings.Builder{}
	destruct := &strings.Builder{}
	code := &strings.Builder{}
	fmt.Fprintf(code, "%s%s(", Namespace, s.Fn)

	for i, par := range s.Args {
		cod, err := c.addNode(par)
		if err != nil {
			return nil, err
		}

		setup.WriteString(cod.Setup)
		destruct.WriteString(cod.Destruct)
		code.WriteString(cod.Code)
		if i != len(s.Args)-1 {
			if len(cod.Setup) > 0 {
				setup.WriteString(";\n")
			}
			if len(cod.Destruct) > 0 {
				destruct.WriteString(";\n")
			}
			code.WriteString(", ")
		}
	}

	code.WriteString(")")

	return &Value{
		Setup:    setup.String(),
		Destruct: destruct.String(),
		Code:     code.String(),
		Type:     s.Type(),
	}, nil
}

func (c *CGen) addReturn(s *ir.ReturnStmt) (*Value, error) {
	v, err := c.addNode(s.Value)
	if err != nil {
		return nil, err
	}

	return &Value{
		Setup:    v.Setup,
		Destruct: v.Destruct,
		Code:     fmt.Sprintf("return %s", v.Code),
		Type:     types.NULL,
	}, nil
}
