package cgen

import (
	"fmt"
	"strconv"

	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/tokenizer"
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

	case types.BOOL:
		code := "true"
		if !val.Value.(bool) {
			code = "false"
		}
		return &Value{
			Code: code,
			Type: types.INT,
		}, nil

	case types.STRING:
		c.RequireSnippet("strings.c")
		varname := fmt.Sprintf("%sstring_%d", Namespace, c.tmpcnt)
		c.tmpcnt++
		setup := fmt.Sprintf("string* %s = string_new(\"%s\", %d);", varname, val.Value.(string), len(val.Value.(string)))
		c.scope.AddFree(fmt.Sprintf("string_free(%s);", varname))

		return &Value{
			Setup:   setup,
			Code:    varname,
			Type:    types.STRING,
			CanGrab: true,
			Grab:    fmt.Sprintf("%s->refs++;", varname),
		}, nil

	default:
		return nil, val.Pos().Error("unknown const type: %s", val.Type().String())
	}
}

func (c *CGen) addDef(pos *tokenizer.Pos, s *ir.DefCall) (*Value, error) {
	v := c.ir.Variables[s.Variable]
	code := fmt.Sprintf("%s %s%s%d", c.GetCType(s.Typ), Namespace, v.Name, v.ID)
	_, exists := dynamicTyps[s.Typ.BasicType()]
	if exists {
		freeCode := c.GetFreeCode(s.Typ, fmt.Sprintf("%s%s%d", Namespace, v.Name, v.ID))
		c.scope.AddFree(freeCode)
	}

	destruct := ""
	if types.ARRAY.Equal(s.Typ) {
		destruct = fmt.Sprintf("%s%s%d = array_new(sizeof(%s), 1);", Namespace, v.Name, v.ID, c.GetCType(s.Typ.(*types.ArrayType).ElemType))
	}
	if types.MAP.Equal(s.Typ) {
		_, exists := c.mapFns[c.getTypName(s.Typ)]
		var freeFn *string
		if !exists {
			var err error
			mapTyp := s.Typ.(*types.MapType)
			freeFn, err = c.addMapTyp(pos, mapTyp)
			if err != nil {
				return nil, err
			}
		} else {
			frFn := "NULL"
			mapTyp := s.Typ.(*types.MapType)
			_, exists1 := dynamicTyps[mapTyp.KeyType.BasicType()]
			_, exists2 := dynamicTyps[mapTyp.ValType.BasicType()]
			if exists1 || exists2 {
				frFn = fmt.Sprintf("%s_free", c.getTypName(mapTyp))
			}
			freeFn = &frFn
		}
		destruct = fmt.Sprintf("%s%s%d = map_new(sizeof(struct %s), %s_compare, %s_hash, %s);", Namespace, v.Name, v.ID, c.getTypName(s.Typ), c.getTypName(s.Typ), c.getTypName(s.Typ), *freeFn)
	}
	return &Value{
		Code:     code,
		Destruct: destruct,
		Type:     types.NULL,
	}, nil
}

func (c *CGen) addVarExpr(s *ir.VariableExpr) (*Value, error) {
	v := c.ir.Variables[s.Variable]
	name := Namespace + v.Name + strconv.Itoa(v.ID)
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

func (c *CGen) addAssign(s *ir.AssignStmt) (*Value, error) {
	// Is map assign?
	a, ok := s.Variable.(*ir.IndexExpr)
	if ok && types.MAP.Equal(a.Value.Type()) {
		m, err := c.addNode(a.Value)
		if err != nil {
			return nil, err
		}
		k, err := c.addNode(a.Index)
		if err != nil {
			return nil, err
		}
		v, err := c.addNode(s.Value)
		if err != nil {
			return nil, err
		}
		return c.addMapSet(m, k, v, a.Value.Type().(*types.MapType))
	}

	lhs, err := c.addNode(s.Variable)
	if err != nil {
		return nil, err
	}

	v, err := c.addNode(s.Value)
	if err != nil {
		return nil, err
	}
	setup := ""
	grabCode := ""
	_, exists := dynamicTyps[s.Type().BasicType()]
	if exists {
		setup = c.GetFreeCode(s.Type(), lhs.Code) + ";\n" + v.Grab
		grabCode = c.GetGrabCode(s.Type(), lhs.Code)
	}
	code := fmt.Sprintf("%s = %s", lhs.Code, v.Code)
	if v.CanGrab {
		setup = JoinCode(setup, v.Grab)
	}

	return &Value{
		Setup:    JoinCode(lhs.Setup, v.Setup, setup),
		Destruct: JoinCode(lhs.Destruct, v.Destruct),
		Code:     code,
		Type:     s.Type(),
		CanGrab:  exists,
		Grab:     grabCode,
	}, nil
}

func (c *CGen) addMathExpr(s *ir.MathExpr) (*Value, error) {
	lhs, err := c.addNode(s.Lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := c.addNode(s.Rhs)
	if err != nil {
		return nil, err
	}
	return &Value{
		Setup:    JoinCode(lhs.Setup, rhs.Setup),
		Destruct: JoinCode(lhs.Destruct, rhs.Destruct),
		Code:     fmt.Sprintf("((%s)%s %s (%s)%s)", c.GetCType(s.Type()), lhs.Code, s.Op, c.GetCType(s.Type()), rhs.Code),
		Type:     s.Type(),
	}, nil
}

var opMap = map[tokenizer.Op]string{
	tokenizer.Eq: "==",
	tokenizer.Ne: "!=",
	tokenizer.Lt: "<",
	tokenizer.Gt: ">",
}

func (c *CGen) addComparison(s *ir.ComparisonExpr) (*Value, error) {
	lhs, err := c.addNode(s.Lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := c.addNode(s.Rhs)
	if err != nil {
		return nil, err
	}
	return &Value{
		Setup:    JoinCode(lhs.Setup, rhs.Setup),
		Destruct: JoinCode(lhs.Destruct, rhs.Destruct),
		Code:     fmt.Sprintf("(%s)%s %s (%s)%s", c.GetCType(s.Typ), lhs.Code, opMap[s.Op], c.GetCType(s.Typ), rhs.Code),
		Type:     s.Type(),
	}, nil
}

func (c *CGen) addStringEq(s *ir.StringEq) (*Value, error) {
	lhs, err := c.addNode(s.Lhs)
	if err != nil {
		return nil, err
	}
	rhs, err := c.addNode(s.Rhs)
	if err != nil {
		return nil, err
	}

	return &Value{
		Setup:    JoinCode(lhs.Setup, rhs.Setup),
		Destruct: JoinCode(lhs.Destruct, rhs.Destruct),
		Code:     fmt.Sprintf("string_equal(%s, %s)", lhs.Code, rhs.Code),
		Type:     types.BOOL,
	}, nil
}

func (c *CGen) addStringCast(s *ir.StringCast) (*Value, error) {
	c.RequireSnippet("format.c")

	v, err := c.addNode(s.Arg)
	if err != nil {
		return nil, err
	}
	var code string
	switch {
	case v.Type.Equal(types.INT):
		code = fmt.Sprintf("string_itoa(%s)", v.Code)

	case v.Type.Equal(types.FLOAT):
		code = fmt.Sprintf("string_ftoa(%s)", v.Code)

	case v.Type.Equal(types.BOOL):
		code = fmt.Sprintf("string_btoa(%s)", v.Code)
	}

	varname := fmt.Sprintf("%sstring_%d", Namespace, c.tmpcnt)
	c.tmpcnt++
	c.scope.AddFree(fmt.Sprintf("string_free(%s);", varname))

	return &Value{
		Setup:    JoinCode(v.Setup, fmt.Sprintf("string* %s = %s;", varname, code)),
		Destruct: v.Destruct,
		Code:     varname,
		Type:     types.STRING,
		CanGrab:  true,
		Grab:     fmt.Sprintf("%s->refs++;", varname),
	}, nil
}
