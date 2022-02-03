package cgen

import (
	"fmt"
	"strings"
)

func (c *CGen) Build() string {
	out := &strings.Builder{}
	for imp := range c.imports {
		fmt.Fprintf(out, "#include <%s>\n", imp)
	}
	out.WriteString("\n")
	out.WriteString(c.top.String())
	return out.String()
}
