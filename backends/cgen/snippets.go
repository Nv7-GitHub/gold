package cgen

import (
	"bufio"
	"embed"
	"io"
	"path"
	"strings"
)

//go:embed snippets/*.c
var snippetsData embed.FS

var snippets = make(map[string]snippet)

type snippet struct {
	imports []string
	code    string
}

func init() {
	files, err := snippetsData.ReadDir("snippets")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		f, err := snippetsData.Open(path.Join("snippets", file.Name()))
		if err != nil {
			panic(err)
		}

		// Read code
		buf := bufio.NewReader(f)
		imports := make([]string, 0)
		for {
			line, _, err := buf.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			ln := string(line)

			if strings.HasPrefix(ln, "#include <") {
				imp := strings.TrimSuffix(strings.TrimPrefix(ln, "#include <"), ">")
				imports = append(imports, imp)
			} else {
				break
			}
		}

		code, err := io.ReadAll(buf)
		if err != nil {
			panic(err)
		}

		f.Close()

		snippets[file.Name()] = snippet{
			imports: imports,
			code:    string(code),
		}
	}
}

func (c *CGen) RequireSnippet(name string) error {
	_, ok := c.snippets[name]
	if ok {
		return nil // Already added
	}

	// Resolve imports
	snip := snippets[name] // No check, assumes code calling is safe
	for _, imp := range snip.imports {
		c.imports[imp] = empty{}
	}

	// Add code
	c.top.WriteString(snip.code)
	c.top.WriteString("\n")

	c.snippets[name] = empty{}

	return nil
}
