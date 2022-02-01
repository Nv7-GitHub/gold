package main

import (
	_ "embed"
	"fmt"

	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/davecgh/go-spew/spew"
)

//go:embed examples/hello.gold
var code string

func main() {
	// TODO: Test out blocks, exprs
	stream := tokenizer.NewStream("hello.gold", code)
	tok := tokenizer.NewTokenizer(stream)
	tok.Tokenize()

	parse := parser.NewParser(tok)
	err := parse.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	bld := ir.NewBuilder()
	ir, err := bld.Build(parse)
	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Config.DisablePointerAddresses = true
	spew.Dump(ir)
}
