package tokenizer

type Tokenizer struct {
	stream *Stream

	Tokens []Token
	pos    int
}

func (t *Tokenizer) Filename() string {
	return t.stream.codePos.Filename
}

func NewTokenizer(stream *Stream) *Tokenizer {
	return &Tokenizer{
		stream: stream,
		Tokens: make([]Token, 0),
	}
}

func (t *Tokenizer) Eat() bool {
	t.pos++
	return !(t.pos >= len(t.Tokens))
}

func (t *Tokenizer) IsEnd() bool {
	return t.pos >= len(t.Tokens)-1
}

func (t *Tokenizer) CurrTok() Token {
	return t.Tokens[t.pos]
}

func (t *Tokenizer) CurrPos() *Pos {
	ps := t.Tokens[t.pos-1].Pos.Dup()
	ps.Col += len(t.Tokens[t.pos-1].Value)
	return ps
}
