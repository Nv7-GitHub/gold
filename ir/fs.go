package ir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type FS interface {
	Parse(src string) (*parser.Parser, error)
}

type defaultFS struct {
	code map[string]string
}

func (f *defaultFS) Parse(src string) (*parser.Parser, error) {
	code, exists := f.code[src]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", src)
	}
	stream := tokenizer.NewTokenizer(tokenizer.NewStream(src, code))
	stream.Tokenize()
	parser := parser.NewParser(stream)

	err := parser.Parse()
	return parser, err
}

// NewFS makes an FS from a map of filename to string
func NewFS(code map[string]string) FS {
	return &defaultFS{
		code: code,
	}
}

type dirFS struct {
	dir string
}

func (d *dirFS) Parse(name string) (*parser.Parser, error) {
	src, err := os.ReadFile(filepath.Join(d.dir, name))
	if err != nil {
		return nil, err
	}
	stream := tokenizer.NewTokenizer(tokenizer.NewStream(name, string(src)))
	stream.Tokenize()
	parser := parser.NewParser(stream)

	err = parser.Parse()
	return parser, err
}

// NewDirFS makes a new FS based on an existing directory on the filesystem
func NewDirFS(dir string) FS {
	return &dirFS{
		dir: dir,
	}
}

func Build(fs FS, main string) (*IR, error) {
	bld := NewBuilder(fs)
	f, err := fs.Parse(main)
	if err != nil {
		return nil, err
	}
	return bld.Build(f)
}

type ImportCall struct{}

func (i *ImportCall) Type() types.Type { return types.NULL }

func init() {
	builders["import"] = nodeBuilder{
		ParamTyps: []types.Type{types.STRING},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			name, ok := args[0].(*Const)
			if !ok {
				return nil, name.Pos().Error("expected constant for filename")
			}
			if !types.STRING.Equal(name.Type()) {
				return nil, name.Pos().Error("expected string for filename")
			}
			if name.IsIdentifier {
				return nil, name.Pos().Error("expected string literal for filename")
			}

			filename := name.Value.(string)

			// Check if imported
			_, exists := b.alreadyImported[filename]
			if exists {
				return &ImportCall{}, nil
			}

			// Build code
			f, err := b.fs.Parse(filename)
			if err != nil {
				return nil, err
			}
			bld := NewBuilder(b.fs)
			bld.alreadyImported = b.alreadyImported // Make imports affect this one's imports
			bld.Variables = b.Variables             // Make variables append to this one's variables
			ir, err := bld.Build(f)                 // This will add the file to alreadyImported
			if err != nil {
				return nil, err
			}

			// Merge functions
			for _, fn := range ir.Funcs {
				existing, exists := b.Funcs[fn.Name]
				if exists {
					return nil, fn.Pos().Error("function already defined at %s", existing.Pos())
				}
				b.Funcs[fn.Name] = fn
			}

			// Add code
			b.TopLevel = append(b.TopLevel, ir.Nodes...)

			// Return
			return &ImportCall{}, nil
		},
	}
}
