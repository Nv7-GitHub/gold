package cgen

import (
	"fmt"

	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/types"
)

func (c *CGen) addPrint(n *ir.PrintStmt) (*Value, error) {
	c.RequireSnippet("print.c")
	val, err := c.addNode(n.Arg)
	if err != nil {
		return nil, err
	}
	return &Value{
		Setup:    val.Setup,
		Destruct: val.Destruct,

		Code: fmt.Sprintf(`printf("%%.*s\n", %s->len, %s->data)`, val.Code, val.Code),
		Type: types.NULL,
	}, nil
}

func (c *CGen) addExit(n *ir.ExitStmt) (*Value, error) {
	stat, err := c.addNode(n.Status)
	if err != nil {
		return nil, err
	}
	code := c.scope.Code()
	return &Value{
		Setup:    JoinCode(stat.Setup, code),
		Destruct: stat.Destruct,
		Code:     fmt.Sprintf("exit(%s)", stat.Code),
		Type:     types.NULL,
	}, nil
}
