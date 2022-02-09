package cgen

import (
	"fmt"

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
			Grab:    fmt.Sprintf("%s->refs++", varname),
		}, nil

	default:
		return nil, val.Pos().Error("unknown const type: %s", val.Type().String())
	}
}

func (c *CGen) addDef(s *ir.DefCall) (*Value, error) {
	code := fmt.Sprintf("%s %s%s", c.GetCType(s.Typ), Namespace, s.Name)
	_, exists := dynamicTyps[s.Typ]
	if exists {
		freeCode := c.GetFreeCode(s.Typ, Namespace+s.Name)
		c.scope.AddFree(freeCode)
	}
	return &Value{
		Code: code,
		Type: types.NULL,
	}, nil
}

func (c *CGen) addVarExpr(s *ir.VariableExpr) (*Value, error) {
	return &Value{
		Code: fmt.Sprintf("%s%s", Namespace, s.Name),
		Type: s.Type(),
	}, nil
}

func (c *CGen) addAssign(s *ir.AssignStmt) (*Value, error) {
	name := Namespace + s.Variable
	v, err := c.addNode(s.Value)
	if err != nil {
		return nil, err
	}
	setup := v.Setup
	grabCode := ""
	_, exists := dynamicTyps[s.Type()]
	if exists {
		if setup != "" {
			setup += ";"
		}
		setup = c.GetFreeCode(s.Type(), s.Variable) + ";\n" + v.Grab
		grabCode = c.GetGrabCode(s.Type(), name)
	}
	code := fmt.Sprintf("%s = %s", name, v.Code)

	return &Value{
		Setup:   setup,
		Code:    code,
		Type:    s.Type(),
		CanGrab: exists,
		Grab:    grabCode,
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
		Code:     fmt.Sprintf("(%s)%s %s (%s)%s", c.GetCType(s.Type()), lhs.Code, s.Op, c.GetCType(s.Type()), rhs.Code),
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
		Grab:     fmt.Sprintf("%s->refs++", varname),
	}, nil
}
