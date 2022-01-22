package tokenizer

import (
	"unicode"
)

func (t *Tokenizer) Tokenize() {
	for t.stream.HasNext() {
		c := t.stream.Peek(0)
		switch c {
		case '"':
			t.Tokens = append(t.Tokens, t.stringLiteral())

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			t.Tokens = append(t.Tokens, t.numLiteral())

		case rune(Add), rune(Sub), rune(Mul), rune(Div), rune(Mod), rune(Eq), rune(Lt), rune(Gt), rune(And), rune(Or):
			if c == '=' && t.stream.Peek(1) == '>' { // Arrow check
				t.Tokens = append(t.Tokens, Token{
					Type:  Operation,
					Value: "=>",
					Pos:   t.stream.CodePos(),
				})
				t.stream.Eat(1)
			} else {
				t.Tokens = append(t.Tokens, Token{
					Type:  Operator,
					Value: string(c),
					Pos:   t.stream.CodePos(),
				})
			}

		case rune(LParen), rune(RParen), rune(LBrack), rune(RBrack):
			t.Tokens = append(t.Tokens, Token{
				Type:  Operator,
				Value: string(c),
				Pos:   t.stream.CodePos(),
			})

		case '#':
			t.eatComment()

		default:
			if unicode.IsLetter(c) {
				t.Tokens = append(t.Tokens, t.identifier())
			} else {
				t.stream.Eat(1)
			}
		}
	}
}

func (t *Tokenizer) stringLiteral() Token {
	pos := t.stream.CodePos()
	t.stream.Eat(1)
	val := ""
	for t.stream.HasNext() {
		if t.stream.Peek(0) == '\\' {
			switch t.stream.Peek(1) {
			case 'n':
				val += "\n"

			case '\\':
				val += "\\"
			}

			t.stream.Eat(1)
		} else {
			val += string(t.stream.Peek(0))
		}

		if t.stream.Peek(1) == '"' {
			t.stream.Eat(2)
			break
		}

		t.stream.Eat(1)
	}
	return Token{
		Type:  StringLiteral,
		Value: val,
		Pos:   pos,
	}
}

func (t *Tokenizer) numLiteral() Token {
	pos := t.stream.CodePos()
	val := ""
	for t.stream.HasNext() {
		val += string(t.stream.Peek(0))

		c := t.stream.Peek(1)
		if c != '.' && !unicode.IsDigit(c) {
			break
		}

		t.stream.Eat(1)
	}
	return Token{
		Type:  NumberLiteral,
		Value: val,
		Pos:   pos,
	}
}

func (t *Tokenizer) identifier() Token {
	pos := t.stream.CodePos()
	val := ""
	for t.stream.HasNext() {
		val += string(t.stream.Peek(0))

		c := t.stream.Peek(1)
		if !unicode.IsLetter(c) {
			t.stream.Eat(1)
			break
		}

		t.stream.Eat(1)
	}
	return Token{
		Type:  Identifier,
		Value: val,
		Pos:   pos,
	}
}

func (t *Tokenizer) eatComment() {
	for t.stream.Peek(0) != '\n' {
		t.stream.Eat(1)
	}
}
