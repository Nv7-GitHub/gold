package main

import (
	_ "embed"
	"fmt"

	"github.com/Nv7-Github/gold/tokenizer"
)

//go:embed examples/hello.gold
var code string

func main() {
	stream := tokenizer.NewStream("hello.gold", code)
	tok := tokenizer.NewTokenizer(stream)
	tok.Tokenize()

	fmt.Println(tok.Tokens)
}
