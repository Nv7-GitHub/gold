package cgen

import (
	"strings"

	"github.com/Nv7-Github/gold/ir"
)

func Indent(code string) string {
	vals := strings.Split(code, "\n")
	out := &strings.Builder{}
	for i, line := range vals {
		out.WriteString("\t")
		out.WriteString(line)
		if i != len(vals)-1 {
			out.WriteString("\n")
		}
	}
	return out.String()
}

func JoinCode(vals ...string) string {
	out := &strings.Builder{}
	for i, val := range vals {
		if len(val) > 0 {
			out.WriteString(val)
			if i != len(vals)-1 {
				out.WriteString("\n")
			}
		}
	}
	return out.String()
}

func (c *CGen) BuildStmts(stmts []ir.Node) (string, error) {
	code := &strings.Builder{}
	c.scope.Push()
	for i, node := range stmts {
		v, err := c.addNode(node)
		if err != nil {
			return "", err
		}
		if len(v.Setup) > 0 {
			code.WriteString(v.Setup)
			code.WriteString("\n")
		}
		code.WriteString(v.Code)
		code.WriteString(";")
		if i != len(stmts)-1 {
			code.WriteString("\n")
		}
		if len(v.Destruct) > 0 {
			code.WriteString(v.Destruct)
			code.WriteString("\n")
		}
	}
	cleanup := c.scope.Pop()
	if len(cleanup) > 0 {
		code.WriteString("\n")
		code.WriteString(cleanup)
	}
	return code.String(), nil
}
