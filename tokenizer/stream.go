package tokenizer

type Stream struct {
	code    []rune
	pos     int
	codePos *Pos
}

func (s *Stream) Peek(off int) rune {
	if len(s.code) <= s.pos+off {
		return rune(0)
	}
	return s.code[s.pos+off]
}

func (s *Stream) Eat(amount int) {
	for i := 0; i < amount; i++ {
		switch s.code[s.pos+i] {
		case '\n':
			s.codePos.NextLine()

		default:
			s.codePos.NextCol()
		}
	}

	s.pos += amount
}

func (s *Stream) CodePos() *Pos {
	return s.codePos.Dup()
}

func (s *Stream) HasNext() bool {
	return s.pos <= len(s.code)-1
}

func NewStream(filename string, code string) *Stream {
	return &Stream{code: []rune(code), codePos: NewPos(filename)}
}
