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

	parse := parser.NewParser(tok)
	err := parse.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, node := range parse.Nodes {
		fmt.Println(node)
	}
}
