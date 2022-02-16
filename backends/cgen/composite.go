package cgen

import (
	"fmt"

	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/types"
)

func (c *CGen) addAppendStmt(s *ir.AppendStmt) (*Value, error) {
	c.RequireSnippet("arrays.c")

	array, err := c.addNode(s.Array)
	if err != nil {
		return nil, err
	}

	v, err := c.addNode(s.Val)
	if err != nil {
		return nil, err
	}

	setup := ""
	if v.CanGrab {
		setup = v.Grab
	}

	return &Value{
		Setup:    JoinCode(array.Setup, v.Setup, setup),
		Destruct: JoinCode(array.Destruct, v.Destruct),
		Code:     fmt.Sprintf("array_append(%s, %s)", array.Code, v.Code),
		Type:     types.NULL,
	}, nil
}

func (c *CGen) addIndexExpr(s *ir.IndexExpr) (*Value, error) {
	v, err := c.addNode(s.Value)
	if err != nil {
		return nil, err
	}
	ind, err := c.addNode(s.Index)
	if err != nil {
		return nil, err
	}

	// Assumes arrays, TODO: support maps
	name := fmt.Sprintf("((%s)array_ind(%s, %s))", c.GetCType(s.Type()), v.Code, ind.Code)
	grabCode := ""
	_, exists := dynamicTyps[s.Type().BasicType()]
	if exists {
		grabCode = c.GetGrabCode(s.Type(), name)
	}

	return &Value{
		Code: name,
		Type: s.Type(),

		CanGrab: exists,
		Grab:    grabCode,
	}, nil
}
