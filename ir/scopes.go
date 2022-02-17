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
	s.scopes = s.scopes[:len(s.scopes)-1]
}

func (s *Scope) AddVar(v *Variable) {
	s.Variables[v.Name] = v
}

func (s *ScopeStack) GetVar(name string) (*Variable, bool) {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		v, exists := s.scopes[i].Variables[name]
		if exists {
			return v, true
		}
	}
	return nil, false
}
