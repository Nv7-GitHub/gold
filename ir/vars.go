package ir

import (
	"github.com/Nv7-Github/gold/parser"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

type DefCall struct {
	Name string
	Typ  types.Type
}

func (d *DefCall) Type() types.Type { return types.NULL }

func init() {
	builders["def"] = nodeBuilder{
		ParamTyps: []types.Type{types.STRING, types.STRING},
		Build: func(b *Builder, pos *tokenizer.Pos, args []Node) (Call, error) {
			// Get name
			name, ok := args[0].(*Const)
			if !ok {
				return nil, args[0].Pos().Error("variable name must be constant")
			}
			if !name.IsIdentifier {
				return nil, args[0].Pos().Error("variable name must be identifier")
			}

			// Get type
			typ, ok := args[1].(*Const)
			if !ok {
				return nil, args[1].Pos().Error("variable type must be constant")
			}
			if !typ.IsIdentifier {
				return nil, args[1].Pos().Error("variable type must be identifier")
			}

			// Parse type
			t, err := types.ParseType(typ.Value.(string))
			if err != nil {
				return nil, args[1].Pos().Error("%s", err.Error())
			}

			// Save
			if b.Scope.Curr().AddVar(name.Value.(string), &Variable{
				Name: name.Value.(string),
				Type: t,
			}) {
				return nil, pos.Error("variable %s already defined", name.Value.(string))
			}
			return &DefCall{
				Name: name.Value.(string),
				Typ:  t,
			}, nil
		},
	}
}

type AssignStmt struct {
	pos *tokenizer.Pos

	Value    Node
	Variable Node
}

func (a *AssignStmt) Type() types.Type {
	return types.NULL
}

func (a *AssignStmt) Pos() *tokenizer.Pos {
	return a.pos
}

func (b *Builder) buildAssignStmt(n *parser.AssignStmt) (Node, error) {
	vr, err := b.buildNode(n.Variable, true)
	if err != nil {
		return nil, err
	}
	_, ok := vr.(*VariableExpr)
	if !ok {
		_, ok := vr.(*IndexExpr)
		if !ok {
			return nil, vr.Pos().Error("cannot assign to node %T", vr)
		}
	}

	v, err := b.buildNode(n.Value)
	if err != nil {
		return nil, err
	}

	if !v.Type().Equal(vr.Type()) {
		return nil, v.Pos().Error("cannot assign %s to %s", v.Type(), vr.Type())
	}

	return &AssignStmt{
		pos:      n.Pos(),
		Value:    v,
		Variable: vr,
	}, nil
}
