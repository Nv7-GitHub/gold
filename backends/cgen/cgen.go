package cgen

import (
	"strings"

	"github.com/Nv7-Github/gold/ir"
)

type empty struct{}

type CGen struct {
	ir *ir.IR

	top      *strings.Builder
	snippets map[string]empty
	imports  map[string]empty
}

func NewCGen(i *ir.IR) *CGen {
	return &CGen{
		ir: i,

		top:      &strings.Builder{},
		snippets: make(map[string]empty),
		imports:  make(map[string]empty),
	}
}
