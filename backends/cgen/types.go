package cgen

import (
	"github.com/Nv7-Github/gold/types"
)

const Namespace = "gold__"

type Value struct {
	Setup    string
	Destruct string

	Code string
	Type types.Type
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
		return "array*"

	case types.MAP:
		return "map_not_implemented"

	default:
		return "unknown"
	}
}
