package cgen

import (
	"strings"
)

type scope struct {
	toFree []string
}

type Stack struct {
	scopes []scope
}

func NewStack() *Stack {
	return &Stack{
		scopes: make([]scope, 0),
	}
}

func (s *Stack) Push() {
	s.scopes = append(s.scopes, scope{toFree: make([]string, 0)})
}

func (s *Stack) Pop() string {
	code := &strings.Builder{}
	sc := s.scopes[len(s.scopes)-1]
	for i, line := range sc.toFree {
		code.WriteString(line)
		if i != len(sc.toFree)-1 {
			code.WriteString("\n")
		}
	}
	s.scopes = s.scopes[:len(s.scopes)-1]
	return code.String()
}

func (s *Stack) AddFree(code string) {
	s.scopes[len(s.scopes)-1].toFree = append(s.scopes[len(s.scopes)-1].toFree, code)
}
