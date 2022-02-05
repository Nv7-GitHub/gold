package cgen

import "strings"

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
