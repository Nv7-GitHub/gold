package ir

func (s *ScopeStack) PushScope(sc *Scope) {
	s.scopes = append(s.scopes, sc)
	s.scopecnt[sc.Type]++
	sc.parent = s
}

func (s *ScopeStack) Curr() *Scope {
	return s.scopes[len(s.scopes)-1]
}

func (s *ScopeStack) Pop() {
	c := s.Curr()
	c.parent = nil
	s.scopecnt[c.Type]--
	for v := range c.Variables {
		delete(s.vars, v)
	}
}

func (s *Scope) AddVar(name string, v *Variable) bool {
	_, exists := s.parent.vars[name]
	if exists {
		return true
	}
	s.Variables[name] = v
	s.parent.vars[name] = v
	return false
}

func (s *ScopeStack) GetVar(name string) (*Variable, bool) {
	v, exists := s.vars[name]
	return v, exists
}
