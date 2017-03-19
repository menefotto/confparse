package confparse

import (
	"io"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	s   *Lexer
	buf struct {
		tok    Token
		values []string
		n      int
	}
}

func newParser(r io.Reader) *Parser {
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

//parser does not take into consideration whitespaces ever
func (p *Parser) Parse() (item *itemType) {
	item = p.scan()
	if item.Tok == WHITESPACE {
		item = p.scan()
	}
	return
}

type config struct {
	C map[string]map[string]string
}

func newConfig() *config {
	conf := &config{C: make(map[string]map[string]string, 0)}
	conf.C["default"] = make(map[string]string, 0)
	conf.C["default"]["version"] = "0.1"
	return conf

}

func (c *config) getValue(section, key string, i *iniParser) (string, error) {
	sec, ok := c.C[section]
	if !ok {
		return "", NewParserError(SEC_NOT_FOUND.Error(), section, key,
			i.errorLine(key))
	}
	val, ok := sec[key]
	if !ok {
		return "", NewParserError(KEY_NOT_FOUND.Error(), section, key,
			i.errorLine(key))
	}

	return val, nil

}

type iniParser struct {
	p *Parser
	c *config
}

func NewParserFromFile(confname string) (*iniParser, error) {
	f, err := os.Open(confname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p := NewParser(f)
	p.Parse()

	return p, nil
}

func NewParser(r io.Reader) *iniParser {
	return &iniParser{p: newParser(r), c: newConfig()}
}

func (i *iniParser) Parse() {
	var lastsection string

	for {
		item := i.p.Parse()

		switch {
		case item.Tok == EOF:
			return
		case item.Tok == KEY_VALUE:
			i.c.C[lastsection][item.Values[0]] = item.Values[1]
		case item.Tok == SECTION:
			lastsection = item.Values[0]
			i.c.C[item.Values[0]] = make(map[string]string, 0)

		}
	}
}

func (i *iniParser) GetBool(section, key string) (bool, error) {
	value, err := i.c.getValue(section, key, i)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, NewParserError(err.Error(), section, key, i.errorLine(key))
	}

	return b, nil

}

func (i *iniParser) GetInt(section, key string) (int64, error) {
	value, err := i.c.getValue(section, key, i)
	if err != nil {
		return -1, err
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1, NewParserError(err.Error(), section, key, i.errorLine(key))
	}

	return n, nil

}

func (i *iniParser) GetFloat(section, key string) (float64, error) {
	value, err := i.c.getValue(section, key, i)
	if err != nil {
		return -0.1, err
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return -1, NewParserError(err.Error(), section, key, i.errorLine(key))

	}

	return f, nil

}

func (i *iniParser) GetString(section, key string) (string, error) {
	value, err := i.c.getValue(section, key, i)
	if err != nil {
		return " ", err
	}
	return value, nil
}

func (i *iniParser) GetSlice(section, key string) ([]string, error) {
	value, err := i.c.getValue(section, key, i)
	if err != nil {
		return []string{""}, NewParserError(err.Error(), section, key, i.errorLine(key))
	}

	return strings.Split(value, ","), nil

}

func (i *iniParser) errorLine(word string) int {
	lineno, err := i.p.s.findLine(word)
	if err == io.EOF {
		return lineno
	}
	if err != nil {
		return -1
	}
	return lineno

}
