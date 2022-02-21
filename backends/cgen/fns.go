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
	fn := c.ir.Funcs[s.Fn]
	tmp := c.tmpcnt
	c.tmpcnt++
	varname := fmt.Sprintf("call_ret%d", tmp)
	fmt.Fprintf(code, "%s %s = %s%s(", c.GetCType(fn.RetType), varname, Namespace, s.Fn)

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

	_, exists := dynamicTyps[fn.RetType.BasicType()]
	if exists {
		c.scope.AddFree(c.GetFreeCode(fn.RetType, varname))
	}

	return &Value{
		Setup:    JoinCode(setup.String(), code.String()),
		Destruct: destruct.String(),
		Code:     varname,
		Type:     s.Type(),
	}, nil
}

func (c *CGen) addReturn(s *ir.ReturnStmt) (*Value, error) {
	v, err := c.addNode(s.Value)
	if err != nil {
		return nil, err
	}
	destruct := v.Destruct
	if v.CanGrab {
		destruct = JoinCode(destruct, v.Grab)
	}

	return &Value{
		Setup:    v.Setup,
		Destruct: destruct,
		Code:     fmt.Sprintf("return %s", v.Code),
		Type:     types.NULL,
	}, nil
}
