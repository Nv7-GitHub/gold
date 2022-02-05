package ir

import (
	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

func (b *Builder) getFunctionDef(p *parser.BlockStmt) (*Func, []parser.Node, error) {
	if p.Fn != "func" || len(p.Args) != 1 {
		return nil, nil, nil
	}

	// Get name
	fnNameV := p.Args[0]
	fnName, ok := fnNameV.(*parser.Const)
	if !ok {
		return nil, nil, fnNameV.Pos().Error("expected constant for function name")
	}
	if !fnName.Type.Equal(types.STRING) {
		return nil, nil, fnNameV.Pos().Error("expected string for function name")
	}
	if !fnName.IsIdentifier {
		return nil, nil, fnNameV.Pos().Error("expected identifier for function name")
	}

	paramTyps := make([]FuncParam, 0)
	for len(p.Stmts) > 0 {
		// Get param, check if param
		s := p.Stmts[0]
		c, ok := s.(*parser.CallStmt)
		if !ok {
			break
		}
		if c.Fn != "param" {
			break
		}
		if len(c.Args) != 2 {
			return nil, nil, c.Pos().Error("expected name, type for declaration \"param\"")
		}

		// Parse param
		name, ok := c.Args[0].(*parser.Const)
		if !ok {
			return nil, nil, name.Pos().Error("expected constant for param name")
		}
		if !name.Type.Equal(types.STRING) {
			return nil, nil, name.Pos().Error("expected string for param name")
		}
		if !name.IsIdentifier {
			return nil, nil, name.Pos().Error("expected identifier for param name")
		}

		typV, ok := c.Args[1].(*parser.Const)
		if !ok {
			return nil, nil, typV.Pos().Error("expected constant for param type")
		}
		if !typV.Type.Equal(types.STRING) {
			return nil, nil, typV.Pos().Error("expected string for param type")
		}
		if !typV.IsIdentifier {
			return nil, nil, typV.Pos().Error("expected identifier for param type")
		}
		typ, err := types.ParseType(typV.Val.(string))
		if err != nil {
			return nil, nil, typV.Pos().Error("%s", err.Error())
		}

		// Add param
		paramTyps = append(paramTyps, FuncParam{
			Name: name.Val.(string),
			Type: typ,
		})
		p.Stmts = p.Stmts[1:]
	}

	// Get return type
	retType := types.Type(types.NULL)
	if len(p.Stmts) == 0 {
		return nil, nil, p.Pos().Error("expected return type")
	}
	s := p.Stmts[0]
	v, ok := s.(*parser.CallStmt)
	if ok && v.Fn == "returns" {
		// Parse return type
		typV, ok := v.Args[0].(*parser.Const)
		if !ok {
			return nil, nil, typV.Pos().Error("expected constant for return type")
		}
		if !typV.Type.Equal(types.STRING) {
			return nil, nil, typV.Pos().Error("expected string for return type")
		}
		if !typV.IsIdentifier {
			return nil, nil, typV.Pos().Error("expected identifier for return type")
		}
		typ, err := types.ParseType(typV.Val.(string))
		if err != nil {
			return nil, nil, typV.Pos().Error("%s", err.Error())
		}
		retType = typ

		p.Stmts = p.Stmts[1:]
	}

	// Return
	return &Func{
		Name:    fnName.Val.(string),
		Params:  paramTyps,
		RetType: retType,
		pos:     p.Pos(),
	}, p.Stmts, nil
}

func (b *Builder) functionPass(p *parser.Parser) error {
	b.Funcs = make(map[string]*Func)
	toBuild := make(map[string][]parser.Node)
	toRemove := make(map[int]struct{})
	for i, stmt := range p.Nodes {
		_, ok := stmt.(*parser.BlockStmt)
		if ok {
			// Is function def?
			f, stmts, err := b.getFunctionDef(stmt.(*parser.BlockStmt))
			if err != nil {
				return err
			}
			if f != nil {
				_, exists := b.Funcs[f.Name]
				if exists {
					continue
				}

				b.Funcs[f.Name] = f
				toBuild[f.Name] = stmts
				toRemove[i] = struct{}{}
			}
		}
	}

	// Remove
	newNodes := make([]parser.Node, 0, len(p.Nodes))
	for i, stmt := range p.Nodes {
		_, exists := toRemove[i]
		if !exists {
			newNodes = append(newNodes, stmt)
		}
	}
	p.Nodes = newNodes

	// Build toBuild
	for name, stmts := range toBuild {
		fn := b.Funcs[name]

		// Prepare scope
		s := NewScope(ScopeTypeFunction)
		s.FuncName = name
		b.Scope.PushScope(s)
		for _, param := range fn.Params {
			s.AddVar(param.Name, &Variable{
				Name: param.Name,
				Type: param.Type,
			})
		}

		// Build
		body := make([]Node, len(stmts))
		var err error
		for i, stmt := range stmts {
			body[i], err = b.buildNode(stmt)
			if err != nil {
				return err
			}
		}
		if b.Scope.Curr().Type != ScopeTypeFunction {
			return fn.Pos().Error("scope not closed, missing \"end\"?")
		}
		b.Scope.Pop()

		fn.Body = body
	}

	return nil
}

type ReturnStmt struct {
	Value Node
}

func (r *ReturnStmt) Type() types.Type { return types.NULL }

type CallStmt struct {
	Fn   string
	Args []Node

	typ types.Type
}

func (c *CallStmt) Type() types.Type { return c.typ }

func init() {
	builders["return"] = nodeBuilder{
		ParamTyps: []types.Type{types.VARIADIC},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			if !b.Scope.HasScope(ScopeTypeFunction) {
				return nil, pos.Error("return statement outside of function")
			}
			var retTyp types.Type
			switch len(args) {
			case 0:
				retTyp = types.NULL

			case 1:
				retTyp = args[0].Type()

			default:
				return nil, pos.Error("expected 0 or 1 arguments")
			}

			fn := b.Funcs[b.Scope.GetScopeByType(ScopeTypeFunction).FuncName]
			if !retTyp.Equal(fn.RetType) {
				return nil, pos.Error("invalid return type: expected %s, got %s", fn.RetType.String(), args[0].Type().String())
			}
			return &ReturnStmt{
				Value: args[0],
			}, nil
		},
	}

	builders["call"] = nodeBuilder{
		ParamTyps: []types.Type{types.STRING, types.VARIADIC},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			// Get function name
			v, ok := args[0].(*Const)
			if !ok {
				return nil, args[0].Pos().Error("expected constant for function name")
			}
			if !v.typ.Equal(types.STRING) {
				return nil, v.Pos().Error("expected string for function name")
			}
			if !v.IsIdentifier {
				return nil, v.Pos().Error("expected identifier for function name")
			}

			// Get function
			fn, exists := b.Funcs[v.Value.(string)]
			if !exists {
				return nil, v.Pos().Error("unknown function: %s", v.Value.(string))
			}

			// Check types
			if len(args)-1 != len(fn.Params) {
				return nil, pos.Error("not enough arguments for function call: expected %d arguments, got %d", len(fn.Params), len(args)-1)
			}
			for i, arg := range args[1:] {
				if !arg.Type().Equal(fn.Params[i].Type) {
					return nil, arg.Pos().Error("invalid argument type: expected %s, got %s", fn.Params[i].Type.String(), arg.Type().String())
				}
			}

			// Make call
			return &CallStmt{
				Fn:   v.Value.(string),
				Args: args[1:],
				typ:  fn.RetType,
			}, nil
		},
	}
}
