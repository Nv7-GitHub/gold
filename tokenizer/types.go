package tokenizer

import "fmt"

type Pos struct {
	Line     int
	Col      int
	Filename string
}

func NewPos(filename string) *Pos {
	return &Pos{Filename: filename}
}

func (p *Pos) NextCol() {
	p.Col++
}

func (p *Pos) NextLine() {
	p.Line++
	p.Col = 0
}

func (p *Pos) Dup() *Pos {
	return &Pos{
		Line:     p.Line,
		Col:      p.Col,
		Filename: p.Filename,
	}
}

func (p *Pos) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line+1, p.Col+1)
}

type TokenType int

const (
	StringLiteral TokenType = iota
	NumberLiteral
	Identifier
	Operator
	Operation
	Parenthesis
	End
)

var tokenNames = map[TokenType]string{
	StringLiteral: "StringLiteral",
	NumberLiteral: "NumberLiteral",
	Identifier:    "Identifier",
	Operator:      "Operator",
	Parenthesis:   "Parenthesis",
	End:           "End",
}

func (t TokenType) String() string {
	return tokenNames[t]
}

type Op rune

const (
	Add Op = '+'
	Sub Op = '-'
	Mul Op = '*'
	Div Op = '/'
	Mod Op = '%'

	Eq  Op = '='
	Ne  Op = '!'
	Lt  Op = '<'
	Gt  Op = '>'
	And Op = '&'
	Or  Op = '|'
)

func (o Op) String() string {
	return string(o)
}

type Paren rune

const (
	LParen Paren = '('
	RParen Paren = ')'
	LBrack Paren = '['
	RBrack Paren = ']'
	LCurly Paren = '{'
	RCurly Paren = '}'
)

type Token struct {
	Type  TokenType
	Value string
	Pos   *Pos
}
