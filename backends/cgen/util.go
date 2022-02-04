package cgen

import "strings"

func Indent(code string) string {
	vals := strings.Split(code, "\n")
	out := &strings.Builder{}
	for _, line := range vals {
		out.WriteString("\t")
		out.WriteString(line)
	}
	return out.String()
}
