package confparse

import (
	"io"
	"os"
	"strconv"
	"strings"
)

type iniParser struct {
	p *parser
	c *config
}

// NewFromFile creates and parse a new configuration from a file name
// returns a valid parsed object and a nil, or an error and nil object
// in the successful case the values are ready to be retrieved
func NewFromFile(confname string) (*iniParser, error) {
	f, err := os.Open(confname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p := New(f)
	p.Parse()

	return p, nil
}

// New creates a new parser from an io.Reader and returns a valid ini parser
// note that the object isn't hasn't yet parsed the configuration Parse has
// to be explicitly called
func New(r io.Reader) *iniParser {
	return &iniParser{p: newParser(r), c: newConfig()}
}

// Parse actually parses the object content, note the object is always in a
// valid state, must be called if the Parser has been created with New, in
// case it has been created with NewFromFile it has already been parsed.
func (i *iniParser) Parse() {
	var lastsection string

	for {
		item := i.p.Scan()

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

// GetBool retrieves a bool value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
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

// GetInt retrieves a int64 value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
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

// GetFloat retrieves a float64 value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
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

// GetString retrieves a string value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
func (i *iniParser) GetString(section, key string) (string, error) {
	value, err := i.c.getValue(section, key, i)
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetSlice retrieves a slice value from named section with key name, returns
// either an error and an invalid value or a nil and a valid value.
func (i *iniParser) GetSlice(section, key string) ([]string, error) {
	value, err := i.c.getValue(section, key, i)
	if err != nil {
		return nil, NewParserError(err.Error(), section, key, i.errorLine(key))
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

type parser struct {
	s   *Lexer
	buf struct {
		tok    Token
		values []string
		n      int
	}
}

func newParser(r io.Reader) *parser {
	return &parser{s: NewLexer(r)}
}

func (p *parser) scan() (item *itemType) {
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

func (p *parser) unscan() { p.buf.n = 1 }

//Scan does not take into consideration white spaces ever nor should be
// called directly.
func (p *parser) Scan() (item *itemType) {
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
