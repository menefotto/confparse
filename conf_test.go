package confparse

import (
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
	parser.Parse()
	val, err := parser.GetString("repos", "base")
	if err != nil {
		t.Log("Error :", err)
	} else {
		t.Logf("Value of key %s is: %s\n", "base", val)
	}

	num, err := parser.GetFloat("repos", "multi")
	if err != nil {
		t.Log("Error :", err)
	} else {
		t.Logf("Value of key %s is: %s\n", "base", num)
	}

	ip, err := parser.GetString("local", "ip")
	if err != nil {
		t.Log("Error :", err)
	} else {
		t.Logf("Value of key %s is: %s\n", "ip", ip)
	}

	pkgs, err := parser.GetSlice("local", "locked_pkgs")
	if err != nil {
		t.Log("Error :", err)
	} else {
		t.Logf("Values of key %-v is: %s\n", "locked pkgs", pkgs)
	}
}

func TestLexer(t *testing.T) {
	f, err := os.Open("sonic.conf")
	defer f.Close()
	if err != nil {
		t.Error(err)
	}
	lexer := NewLexer(f)
	num, err := lexer.findLine("multi")
	if err != nil {
		return
	}
	t.Log("Line number is: ", num)
}
