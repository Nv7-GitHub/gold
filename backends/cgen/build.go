package cgen

import (
	"fmt"
	"strconv"
	"strings"
)

func (c *CGen) Build() (string, error) {
	body := &strings.Builder{}

	// Add function definitions
	for _, fn := range c.ir.Funcs {
		// Build func starter
		body.WriteString(c.GetCType(fn.RetType))
		body.WriteString(" ")
		body.WriteString(Namespace)
		body.WriteString(fn.Name)
		body.WriteString("(")
		for i, par := range fn.Params {
			body.WriteString(c.GetCType(par.Type))
			body.WriteString(" ")
			body.WriteString(Namespace)
			body.WriteString(par.Name)
			body.WriteString(strconv.Itoa(par.VarID))
			if i != len(fn.Params)-1 {
				body.WriteString(", ")
			}
		}
		body.WriteString(");\n")
	}
	body.WriteString("\n")

	// Add functions
	for _, fn := range c.ir.Funcs {
		// Build func starter
		body.WriteString(c.GetCType(fn.RetType))
		body.WriteString(" ")
		body.WriteString(Namespace)
		body.WriteString(fn.Name)
		body.WriteString("(")
		for i, par := range fn.Params {
			body.WriteString(c.GetCType(par.Type))
			body.WriteString(" ")
			body.WriteString(Namespace)
			body.WriteString(par.Name)
			body.WriteString(strconv.Itoa(par.VarID))
			if i != len(fn.Params)-1 {
				body.WriteString(", ")
			}
		}
		body.WriteString(") {\n")

		// Add body
		bdy, err := c.BuildStmts(fn.Body)
		if err != nil {
			return "", err
		}
		body.WriteString(Indent(bdy))

		body.WriteString("\n}\n\n")
	}

	// Build global code
	body.WriteString("int main() {\n")
	bdy, err := c.BuildStmts(c.ir.Nodes)
	if err != nil {
		return "", err
	}
	body.WriteString(Indent(bdy))
	body.WriteString("\n\treturn 0;\n}\n")

	// Build imports & top
	out := &strings.Builder{}
	for imp := range c.imports {
		fmt.Fprintf(out, "#include <%s>\n", imp)
	}
	out.WriteString("\n")
	out.WriteString(c.top.String())
	out.WriteString("\n")
	out.WriteString(c.types.String())
	out.WriteString("\n")
	out.WriteString(body.String())
	return out.String(), nil
}
