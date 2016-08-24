package confparse

import "io"

type Parser struct {
	s   *Lexer
	buf struct {
		tok    Token
		values []string
		n      int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewLexer(r)}
}

func (p *Parser) scan() (item *itemType) {
	// If we have a token on the buffer, then return it.
	if p.buf.values == nil {
		p.buf.values = make([]string, 0)
	}

	if p.buf.n != 0 {
		p.buf.n = 0
		item.Tok = p.buf.tok
		item.Values = append(item.Values, p.buf.values...)
	}

	// Otherwise read the next token from the scanner.
	item = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok = item.Tok
	p.buf.values = append(p.buf.values, item.Values...)

	return
}

func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) Parse() (item *itemType) {
	item = p.scan()
	if item.Tok == WHITESPACE {
		item = p.scan()
	}
	return
}
