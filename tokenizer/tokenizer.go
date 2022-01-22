package tokenizer

type Tokenizer struct {
	stream *Stream

	Tokens []Token
}

func NewTokenizer(stream *Stream) *Tokenizer {
	return &Tokenizer{
		stream: stream,
		Tokens: make([]Token, 0),
	}
}
