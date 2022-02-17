package cgen

import (
	"strings"

	"github.com/Nv7-Github/gold/ir"
)

type empty struct{}

type CGen struct {
	ir *ir.IR

	top      *strings.Builder
	types    *strings.Builder
	freeFns  map[string]empty
	snippets map[string]empty
	imports  map[string]empty

	tmpcnt int

	scope *Stack
}

func NewCGen(i *ir.IR) *CGen {
	return &CGen{
		ir: i,

		top:      &strings.Builder{},
		types:    &strings.Builder{},
		freeFns:  make(map[string]empty),
		snippets: make(map[string]empty),
		imports:  make(map[string]empty),
		tmpcnt:   0,

		scope: NewStack(),
	}
}
