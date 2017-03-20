package confparse

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	parser, err := New("sonic.conf")
	if err != nil {
		t.Fatal(err)
	}

	val, err := parser.GetString("repos.base")
	if err != nil {
		t.Log("Error :", err)
	} else {
		t.Logf("Value of key %s is: %s\n", "base", val)
	}

	num, err := parser.GetFloat("repos.multi")
	if err != nil {
		t.Log("Error :", err)
	} else {
		t.Logf("Value of key %s is: %s\n", "base", num)
	}

	ip, err := parser.GetString("local.ip")
	if err != nil {
		t.Log("Error :", err)
	} else {
		t.Logf("Value of key %s is: %s\n", "ip", ip)
	}

	pkgs, err := parser.GetSlice("local.locked_pkgs")
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
	lexer := newLexer(f)
	num, err := lexer.findLine("multi")
	if err != nil {
		return
	}
	t.Log("Line number is: ", num)
}

func TestFileParser(t *testing.T) {
	_, err := New("sonic.conf")
	if err != nil {
		t.Fatal(err)
	}

	_, err = New("nonexistant.conf")
	if err == nil {
		t.Fatal(err)
	}

}

func TestGetIntBool(t *testing.T) {
	ini, err := New("sonic.conf")
	if err != nil {
		t.Fatal(err)
	}

	n, err := ini.GetInt("local.testint")
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Fatal("Int should be 5!")
	}

	b, err := ini.GetBool("local.testbool")
	if err != nil {
		t.Fatal(err)
	}
	if b != false {
		t.Fatal("B should be false!")
	}

}

func TestBreakGetIntBool(t *testing.T) {
	ini, err := New("sonic.conf")
	if err != nil {
		t.Fatal(err)
	}

	_, err = ini.GetInt("local.testbool")
	if err == nil {
		t.Fatal("Should be an error!")
	}

	_, err = ini.GetBool("local.testint")
	if err == nil {
		t.Fatal("Should be an error!")
	}

}

func TestBreakSliceFload(t *testing.T) {
	ini, err := New("sonic.conf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ini.GetSlice("local.testslice")
	if err == nil {
		t.Fatal("Should be an error!")
	}

	_, err = ini.GetFloat("local.testbool")
	if err == nil {
		t.Fatal("Should be an error!")
	}
}
