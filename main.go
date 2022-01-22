package main

import (
	_ "embed"
	"fmt"

	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
)

//go:embed examples/hello.gold
var code string

func main() {
	stream := tokenizer.NewStream("hello.gold", code)
	tok := tokenizer.NewTokenizer(stream)
	tok.Tokenize()

	parser := parser.NewParser(tok)
	err := parser.Parse()
	if err != nil {
		panic(err)
	}

	for _, node := range parser.Nodes {
		fmt.Println(node)
	}
}
