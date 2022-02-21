package cgen

import (
	"fmt"

	"github.com/Nv7-Github/gold/types"
)

const Namespace = "gold__"

var dynamicTyps = map[types.BasicType]empty{
	types.STRING: {},
	types.ARRAY:  {},
	types.MAP:    {},
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
		c.RequireSnippet("map.c")
		return "map*"

	default:
		return "unknown"
	}
}

func (c *CGen) getTypName(typ types.Type) string {
	switch typ.BasicType() {
	case types.BOOL:
		return "bool"

	case types.INT:
		return "int"

	case types.FLOAT:
		return "float"

	case types.STRING:
		return "string"

	case types.ARRAY:
		return "array_" + c.getTypName(typ.(*types.ArrayType).ElemType)

	case types.MAP:
		return "map_" + c.getTypName(typ.(*types.MapType).KeyType) + "_" + c.getTypName(typ.(*types.MapType).ValType)

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
		freeFn := "NULL"
		if exists {
			name := c.getTypName(elType)
			_, exists := c.freeFns[name]
			if !exists {
				elFree := c.GetFreeCode(elType, fmt.Sprintf("*((%s*)array_get(arr, i))", c.GetCType(elType)))
				loop := fmt.Sprintf("for (int i = 0; i < arr->len; i++) {\n\t%s\n}", elFree)
				fmt.Fprintf(c.types, "void array_free_%s(array* arr) {\n%s\n}\n", name, Indent(loop))
				c.freeFns[name] = empty{}
			}
			freeFn = "array_free_" + name
		}
		return fmt.Sprintf("array_free(%s, %s);", varName, freeFn)

	case types.MAP.Equal(typ):
		return fmt.Sprintf("map_free(%s);", varName)

	default:
		return ""
	}
}

func (c *CGen) GetGrabCode(typ types.Type, varName string) string {
	switch typ.BasicType() {
	case types.STRING, types.ARRAY, types.MAP:
		return fmt.Sprintf("%s->refs++;", varName)

	default:
		return ""
	}
}
