package confparse

import (
	"fmt"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	f, err := os.Open("sonic.conf")
	defer f.Close()
	if err != nil {
		t.Error(err)
	}

	parser := NewParser(f)
	for {
		item := parser.Parse()
		switch {
		case item.Tok == EOF:
			return
		case item.Tok == KEY_VALUE:
			fmt.Println("key value: ", item.Values[0], item.Values[1])
		case item.Tok == COMMENT:
			fmt.Println("comment: ", item.Values[0])
		case item.Tok == SECTION:
			fmt.Println("section: ", item.Values[0])
		}
	}
}
