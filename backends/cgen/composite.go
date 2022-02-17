package cgen

import (
	"fmt"

	"github.com/Nv7-Github/gold/ir"
	"github.com/Nv7-Github/gold/tokenizer"
	"github.com/Nv7-Github/gold/types"
)

func (c *CGen) addAppendStmt(s *ir.AppendStmt) (*Value, error) {
	c.RequireSnippet("arrays.c")

	array, err := c.addNode(s.Array)
	if err != nil {
		return nil, err
	}

	v, err := c.addNode(s.Val)
	if err != nil {
		return nil, err
	}

	setup := ""
	_, exists := dynamicTyps[v.Type.BasicType()]
	if !exists { // If not dynamic, can't get pointer if literal
		tmp := c.tmpcnt
		c.tmpcnt++
		setup = fmt.Sprintf("%s  %sapp_%d = %s;\n", c.GetCType(v.Type), Namespace, tmp, v.Code)
		v.Code = fmt.Sprintf("%sapp_%d", Namespace, tmp)
	}
	if v.CanGrab {
		setup = JoinCode(setup, v.Grab)
	}

	return &Value{
		Setup:    JoinCode(array.Setup, v.Setup, setup),
		Destruct: JoinCode(array.Destruct, v.Destruct),
		Code:     fmt.Sprintf("array_append(%s, &(%s))", array.Code, v.Code),
		Type:     types.NULL,
	}, nil
}

func (c *CGen) addIndexExpr(s *ir.IndexExpr) (*Value, error) {
	// Is map?
	if types.MAP.Equal(s.Value.Type()) {
		m, err := c.addNode(s.Value)
		if err != nil {
			return nil, err
		}
		k, err := c.addNode(s.Index)
		if err != nil {
			return nil, err
		}
		return c.addMapGet(m, k, s.Value.Type().(*types.MapType))
	}

	v, err := c.addNode(s.Value)
	if err != nil {
		return nil, err
	}
	ind, err := c.addNode(s.Index)
	if err != nil {
		return nil, err
	}

	// Assumes arrays, TODO: support maps
	name := fmt.Sprintf("(*((%s*)array_get(%s, %s)))", c.GetCType(s.Type()), v.Code, ind.Code)
	grabCode := ""
	_, exists := dynamicTyps[s.Type().BasicType()]
	if exists {
		grabCode = c.GetGrabCode(s.Type(), name)
	}

	return &Value{
		Code: name,
		Type: s.Type(),

		CanGrab: exists,
		Grab:    grabCode,
	}, nil
}

func (c *CGen) addGrowStmt(s *ir.GrowStmt) (*Value, error) {
	v, err := c.addNode(s.Array)
	if err != nil {
		return nil, err
	}
	ind, err := c.addNode(s.Size)
	if err != nil {
		return nil, err
	}

	return &Value{
		Code: fmt.Sprintf("array_grow_gold(%s, %s)", v.Code, ind.Code),
		Type: types.NULL,
	}, nil
}

func (c *CGen) addMapTyp(pos *tokenizer.Pos, mapTyp *types.MapType) (*string, error) {
	freeFn := "NULL"

	// Add struct
	fmt.Fprintf(c.types, "struct %s {\n\t%s key;\n\t%s val;\n};\n\n", c.getTypName(mapTyp), c.GetCType(mapTyp.KeyType), c.GetCType(mapTyp.ValType))

	// Add compare fn
	fmt.Fprintf(c.types, "int %s_compare(const void *va, const void *vb, void *udata) {\n\tconst struct %s *a = va;\n\tconst struct %s *b = vb;\n\t", c.getTypName(mapTyp), c.getTypName(mapTyp), c.getTypName(mapTyp))
	switch mapTyp.KeyType {
	case types.INT, types.FLOAT:
		c.types.WriteString("return (a->key == b->key) ? 0 : 1;")

	case types.STRING:
		c.types.WriteString("return string_equal(a->key, b->key) ? 0 : 1;")

	default:
		return nil, pos.Error("unsupported map key type: %s", mapTyp.KeyType.String())
	}
	c.types.WriteString("\n}\n\n")

	// Add hash fn
	fmt.Fprintf(c.types, "static inline uint64_t %s_hash(const void *item, uint64_t seed0, uint64_t seed1) {\n\tconst struct %s *v = item;\n\t", c.getTypName(mapTyp), c.getTypName(mapTyp))
	switch mapTyp.KeyType {
	case types.INT, types.FLOAT:
		c.types.WriteString("return *(uint64_t*)(v->key);")

	case types.STRING:
		c.types.WriteString("return hashmap_sip(v->key->data, v->key->len, seed0, seed1);")

	default:
		return nil, pos.Error("unsupported map key type: %s", mapTyp.KeyType.String())
	}
	c.types.WriteString("\n}\n\n")

	// Add free fn
	_, exists1 := dynamicTyps[mapTyp.KeyType.BasicType()]
	_, exists2 := dynamicTyps[mapTyp.ValType.BasicType()]
	if exists1 || exists2 {
		freeFn = fmt.Sprintf("%s_free", c.getTypName(mapTyp))
		fmt.Fprintf(c.types, "void %s_free(void *item) {\n\tconst struct %s *v = item;\n\t", c.getTypName(mapTyp), c.getTypName(mapTyp))
		if exists1 {
			c.types.WriteString(c.GetFreeCode(mapTyp.KeyType, "v->key"))
			if exists2 {
				c.types.WriteString("\n\t")
			}
		}
		if exists2 {
			c.types.WriteString(c.GetFreeCode(mapTyp.ValType, "v->val"))
		}
		c.types.WriteString("\n}\n\n")
	}

	c.mapFns[c.getTypName(mapTyp)] = empty{}

	return &freeFn, nil
}

func (c *CGen) addMapSet(m *Value, k, v *Value, mapTyp *types.MapType) (*Value, error) {
	setup := JoinCode(m.Setup, k.Setup, v.Setup)
	if k.CanGrab {
		setup = JoinCode(setup, k.Grab)
	}
	if v.CanGrab {
		setup = JoinCode(setup, v.Grab)
	}

	// Check if has stuff in key already
	_, exists := dynamicTyps[mapTyp.KeyType.BasicType()]
	_, exists2 := dynamicTyps[mapTyp.ValType.BasicType()]
	if exists || exists2 {
		tmp := c.tmpcnt
		c.tmpcnt++
		mapcheck := fmt.Sprintf("%smapcheck_%d", Namespace, tmp)
		check := fmt.Sprintf("struct %s* %s = hashmap_get(%s->map, &(struct %s){ .key=%s });\nif (%s != NULL) {\n\t%s_free(%s);\n}\n", c.getTypName(mapTyp), mapcheck, m.Code, c.getTypName(mapTyp), k.Code, mapcheck, c.getTypName(mapTyp), mapcheck)
		setup = JoinCode(setup, check)
	}

	return &Value{
		Setup:    setup,
		Destruct: JoinCode(m.Destruct, k.Destruct, v.Destruct),
		Type:     types.NULL,
		Code:     fmt.Sprintf("hashmap_set(%s->map, &(struct %s){ .key=%s, .val=%s })", m.Code, c.getTypName(mapTyp), k.Code, v.Code),
	}, nil
}

func (c *CGen) addMapGet(m *Value, k *Value, mapTyp *types.MapType) (*Value, error) {
	tmp := c.tmpcnt
	c.tmpcnt++
	valTyp := mapTyp.ValType
	name := fmt.Sprintf("%smap_%d", Namespace, tmp)
	setup := JoinCode(m.Setup, k.Setup, fmt.Sprintf("%s %s = (%s)(((struct %s*)hashmap_get(%s->map, &(struct %s){ .key=%s }))->val);", c.GetCType(valTyp), name, c.GetCType(valTyp), c.getTypName(mapTyp), m.Code, c.getTypName(mapTyp), k.Code))

	grabCode := ""
	_, exists := dynamicTyps[valTyp.BasicType()]
	if exists {
		grabCode = c.GetGrabCode(valTyp, name)
	}

	return &Value{
		Setup:    setup,
		Destruct: JoinCode(m.Destruct, k.Destruct),
		Type:     valTyp,
		Code:     name,
		Grab:     grabCode,
		CanGrab:  exists,
	}, nil
}
