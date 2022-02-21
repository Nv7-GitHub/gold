package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Nv7-Github/gold/backends/cgen"
	"github.com/Nv7-Github/gold/ir"
)

//go:embed examples/hello.gold
var code string

func main() {
	ir, err := ir.Build(ir.NewDirFS("examples/calculator"), "calculator.gold")
	if err != nil {
		fmt.Println(err)
		return
	}

	// CGen test
	cgen := cgen.NewCGen(ir)
	cgen.RequireSnippet("strings.c")
	code, err := cgen.Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write & compile
	f, err := os.Create("examples/code.c")
	if err != nil {
		panic(err)
	}
	/*f, err := os.CreateTemp("", "*.c")
	if err != nil {
		panic(err)
	}*/
	_, err = f.WriteString(code)
	if err != nil {
		panic(err)
	}

	currPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("cc", f.Name(), "-o", filepath.Join(currPath, "main"))
	stderr := bytes.NewBuffer(nil)
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return
	}
}
