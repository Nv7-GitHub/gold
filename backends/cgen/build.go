package cgen

import (
	"fmt"
	"strings"
)

func (c *CGen) Build() (string, error) {
	body := &strings.Builder{}

	// Add functions
	for _, fn := range c.ir.Funcs {
		// Build func starter
		body.WriteString(c.GetCType(fn.RetType))
		body.WriteString(" ")
		body.WriteString(fn.Name)
		body.WriteString("(")
		for i, par := range fn.Params {
			body.WriteString(c.GetCType(par.Type))
			body.WriteString(" ")
			body.WriteString(par.Name)
			if i != len(fn.Params)-1 {
				body.WriteString(", ")
			}
		}
		body.WriteString(") {\n")

		// Add body
		c.scope.Push()
		for _, node := range fn.Body {
			v, err := c.addNode(node)
			if err != nil {
				return "", err
			}
			body.WriteString(Indent(v.Setup))
			body.WriteString("\n")
			body.WriteString(Indent(v.Code))
			body.WriteString(";\n")
			body.WriteString(Indent(v.Destruct))
			body.WriteString("\n")
		}
		// Add free code
		body.WriteString(Indent(c.scope.Pop()))

		body.WriteString("}\n\n")
	}

	// Build global code
	body.WriteString("int main() {\n")
	c.scope.Push()
	for _, node := range c.ir.Nodes {
		code, err := c.addNode(node)
		if err != nil {
			return "", err
		}
		body.WriteString(Indent(code.Setup))
		body.WriteString("\n")
		body.WriteString(Indent(code.Code))
		body.WriteString(";\n")
		body.WriteString(Indent(code.Destruct))
		body.WriteString("\n")
	}
	body.WriteString(Indent(c.scope.Pop()))
	body.WriteString("\n\treturn 0;\n}\n")

	// Build imports & top
	out := &strings.Builder{}
	for imp := range c.imports {
		fmt.Fprintf(out, "#include <%s>\n", imp)
	}
	out.WriteString("\n")
	out.WriteString(c.top.String())
	out.WriteString("\n")
	out.WriteString(body.String())
	return out.String(), nil
}
