package cgen

import (
	"fmt"

	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/types"
)

func (c *CGen) addConst(val *ir.Const) (*Value, error) {
	switch val.Type() {
	case types.INT:
		return &Value{
			Code: fmt.Sprintf("%d", val.Value.(int)),
			Type: types.INT,
		}, nil

	case types.FLOAT:
		return &Value{
			Code: fmt.Sprintf("%f", val.Value.(float64)),
			Type: types.INT,
		}, nil

	case types.STRING:
		c.RequireSnippet("strings.c")
		varname := fmt.Sprintf("%sstring_%d", Namespace, c.tmpcnt)
		c.tmpcnt++
		setup := fmt.Sprintf("string* %s = string_new(\"%s\", %d);", varname, val.Value.(string), len(val.Value.(string)))
		c.scope.AddFree(fmt.Sprintf("string_free(%s);", varname))

		return &Value{
			Setup: setup,
			Code:  varname,
			Type:  types.STRING,
		}, nil

	default:
		return nil, val.Pos().Error("unknown const type: %s", val.Type().String())
	}
}
