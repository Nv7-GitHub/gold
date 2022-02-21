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

	_, exists := dynamicTyps[fn.RetType.BasicType()]
	if exists {
		tmp := c.tmpcnt
		c.tmpcnt++
		varname := fmt.Sprintf("call_ret%d", tmp)
		cd := fmt.Sprintf("%s %s = %s;\n", c.GetCType(fn.RetType), varname, code.String())
		c.scope.AddFree(c.GetFreeCode(fn.RetType, varname))

		return &Value{
			Setup:    JoinCode(setup.String(), cd),
			Destruct: destruct.String(),
			Code:     varname,
			Type:     s.Type(),
			CanGrab:  true,
			Grab:     c.GetGrabCode(s.Type(), varname),
		}, nil
	}

	return &Value{
		Setup:    setup.String(),
		Destruct: destruct.String(),
		Code:     code.String(),
		Type:     s.Type(),
	}, nil
}

func (c *CGen) addReturn(s *ir.ReturnStmt) (*Value, error) {
	// Add free code
	setup := c.scope.Code()

	if s.Value == nil { // No args
		return &Value{
			Setup: setup,
			Code:  "return",
			Type:  types.NULL,
		}, nil
	}

	v, err := c.addNode(s.Value)
	if err != nil {
		return nil, err
	}
	_, exists := dynamicTyps[v.Type.BasicType()]
	if exists {
		setup = JoinCode(v.Grab, setup)
	}
	destruct := v.Destruct
	if v.CanGrab {
		destruct = JoinCode(destruct, v.Grab)
	}

	return &Value{
		Setup:    JoinCode(setup, v.Setup),
		Destruct: destruct,
		Code:     fmt.Sprintf("return %s", v.Code),
		Type:     types.NULL,
	}, nil
}
