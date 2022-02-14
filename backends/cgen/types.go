package cgen

import (
	"fmt"

	"github.com/Nv7-Github/gold/types"
)

const Namespace = "gold__"

var dynamicTyps = map[types.BasicType]empty{
	types.STRING: {},
	types.ARRAY:  {},
}

type Value struct {
	Setup    string
	Destruct string

	Code string
	Type types.Type

	CanGrab bool
	Grab    string
}

func (c *CGen) GetCType(typ types.Type) string {
	switch typ.BasicType() {
	case types.BOOL:
		return "bool"

	case types.INT:
		return "int"

	case types.FLOAT:
		return "float"

	case types.NULL:
		return "void"

	case types.STRING:
		c.RequireSnippet("strings.c")
		return "string*"

	case types.ARRAY:
		c.RequireSnippet("arrays.c")
		return "array*"

	case types.MAP:
		return "map_not_implemented"

	default:
		return "unknown"
	}
}

func (c *CGen) GetFreeCode(typ types.Type, varName string) string {
	switch {
	case types.STRING.Equal(typ):
		return fmt.Sprintf("string_free(%s);", varName)

	case types.ARRAY.Equal(typ):
		elType := typ.(*types.ArrayType).ElemType
		_, exists := dynamicTyps[elType.BasicType()]
		if exists {
			tmpV := c.tmpcnt
			c.tmpcnt++
			elFree := c.GetFreeCode(elType, fmt.Sprintf("(%s)%s[i%d]", c.GetCType(elType), varName, tmpV))
			return fmt.Sprintf("for (int i%d = 0; i%d < %s->len; i%d++) {\n\t%s}\narray_free(%s);", tmpV, tmpV, varName, tmpV, elFree, varName)
		}
		return fmt.Sprintf("array_free(%s);", varName)

	default:
		return ""
	}
}

func (c *CGen) GetGrabCode(typ types.Type, varName string) string {
	switch typ {
	case types.STRING:
		return fmt.Sprintf("%s->refs++;", varName)

	default:
		return ""
	}
}
