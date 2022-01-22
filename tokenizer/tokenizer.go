package tokenizer

type Tokenizer struct {
	stream *Stream

	Tokens []Token
	pos    int
}

func NewTokenizer(stream *Stream) *Tokenizer {
	return &Tokenizer{
		stream: stream,
		Tokens: make([]Token, 0),
	}
}

func (t *Tokenizer) Eat() {
	t.pos++
}

func (t *Tokenizer) IsEnd() bool {
	return t.pos >= len(t.Tokens)-1
}

func (t *Tokenizer) CurrTok() Token {
	return t.Tokens[t.pos]
}
