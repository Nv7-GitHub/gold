package cgen

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/types"
)

func (c *CGen) addWhile(s *ir.WhileStmt) (*Value, error) {
	cond, err := c.addNode(s.Cond)
	if err != nil {
		return nil, err
	}

	body, err := c.BuildStmts(s.Body)
	if err != nil {
		return nil, err
	}

	return &Value{
		Setup:    cond.Setup,
		Destruct: cond.Destruct,
		Code:     fmt.Sprintf("while (%s) {\n%s\n}", cond.Code, Indent(body)),
		Type:     types.NULL,
	}, nil
}

func (c *CGen) addIf(s *ir.IfStmt) (*Value, error) {
	cond, err := c.addNode(s.Cond)
	if err != nil {
		return nil, err
	}

	body, err := c.BuildStmts(s.Body)
	if err != nil {
		return nil, err
	}
	elseCode := ""
	if s.Else != nil {
		els, err := c.BuildStmts(s.Else)
		if err != nil {
			return nil, err
		}
		elseCode = fmt.Sprintf(" else {\n%s\n}", Indent(els))
	}

	return &Value{
		Setup:    cond.Setup,
		Destruct: cond.Destruct,
		Code:     fmt.Sprintf("if (%s) {\n%s\n}%s", cond.Code, Indent(body), elseCode),
		Type:     types.NULL,
	}, nil
}

func (c *CGen) getHash(typ types.Type, val *ir.Const) (string, error) {
	switch {
	case types.STRING.Equal(typ):
		v := fnv.New32a()
		v.Write([]byte(val.Value.(string)))
		return fmt.Sprintf("%d", v.Sum32()), nil

	case types.INT.Equal(typ):
		return fmt.Sprintf("%d", val.Value.(int)), nil

	case types.FLOAT.Equal(typ):
		return fmt.Sprintf("%f", val.Value.(float64)), nil

	default:
		return "", fmt.Errorf("cannot switch on type %s", typ)
	}
}

func (c *CGen) getHashCode(typ types.Type, code string) (string, error) {
	switch {
	case types.STRING.Equal(typ):
		c.RequireSnippet("switch.c")
		return fmt.Sprintf("str_hash(%s)", code), nil

	case types.INT.Equal(typ), types.FLOAT.Equal(typ):
		return code, nil

	default:
		return "", fmt.Errorf("cannot switch on type %s", typ)
	}
}

func (c *CGen) addSwitch(s *ir.SwitchStmt) (*Value, error) {
	cond, err := c.addNode(s.Cond)
	if err != nil {
		return nil, err
	}
	out := &strings.Builder{}

	// Get cond
	hashCond, err := c.getHashCode(s.Cond.Type(), cond.Code)
	if err != nil {
		return nil, s.Cond.Pos().Error("%s", err.Error())
	}
	fmt.Fprintf(out, "switch (%s) {\n", hashCond)

	// Make cases
	for _, cs := range s.Cases {
		hash, err := c.getHash(cs.Cond.Type(), cs.Cond)
		if err != nil {
			return nil, cs.Cond.Pos().Error("%s", err.Error())
		}
		fmt.Fprintf(out, "case %s:;\n", hash)

		body, err := c.BuildStmts(cs.Body)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(out, "%s\n\tbreak;\n", Indent(body))
	}

	// Default case
	if s.Default != nil {
		body, err := c.BuildStmts(s.Default.Body)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(out, "default:;\n%s\n\tbreak;\n", Indent(body))
	}
	out.WriteString("}")

	// Make code
	return &Value{
		Setup:    cond.Setup,
		Destruct: cond.Destruct,
		Code:     out.String(),
		Type:     types.NULL,
	}, nil
}
