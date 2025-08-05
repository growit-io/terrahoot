package tool_test

import (
	"os"
	"strings"
	"testing"

	"github.com/growit-io/terrahoot/internal/tool"
)

func TestUsrBinEnv(t *testing.T) {
	// Every POSIX-compliant system should have the env(1) command executable
	// installed in this well-known location.
	tool := tool.New("/usr/bin/env")
	if err := tool.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestGoInPathWithVersionArg_OutputLine(t *testing.T) {
	// Since we're running a test suite written in Go, it is highly likely that
	// the "go" executable can be found in PATH.
	tool := tool.New("go").WithArgs("version")
	line, err := tool.OutputLine()
	if err != nil {
		t.Fatal(err)
	}

	prefix := "go version "
	if !strings.HasPrefix(line, prefix) {
		t.Fatalf("expected output to have prefix %s", prefix)
	}
}

// TestToolWithDir_empty ensures that the tool.Tool.WithDir method works as
// expected when called with the empty string.
func TestToolWithDir_empty(t *testing.T) {
	a, err := tool.New("pwd").WithDir("").OutputLine()
	if err != nil {
		t.Errorf("`pwd` with \"\" as working directory: %v", err)
	}

	b, err := tool.New("pwd").WithDir(".").OutputLine()
	if err != nil {
		t.Errorf("`pwd` with \".\" as working directory: %v", err)
	}

	if a != b {
		t.Error("expected WithDir(\"\") to behave the same as WithDir(\".\")")
	}
}

// TestToolWithDir_dot ensures that the tool.Tool.WithDir method works as
// expected when called with "." as the directory argument.
func TestToolWithDir_dot(t *testing.T) {
	actual, err := tool.New("pwd").WithDir(".").OutputLine()
	if err != nil {
		t.Errorf("`pwd` with \".\" as working directory: %v", err)
	}

	expected, err := os.Getwd()
	if err != nil {
		t.Errorf("os.Getwd: %v", err)
	}

	if expected != actual {
		t.Errorf("expected output of tool.New(\"pwd\").WithDir(\".\") to be the same as the output of `pwd`: %v != %v", actual, expected)
	}
}

// TestToolWithDir_root ensures that the tool.Tool.WithDir method works as
// expected when called with "/" as the directory argument.
func TestToolWithDir_root(t *testing.T) {
	actual, err := tool.New("pwd").WithDir("/").OutputLine()
	if err != nil {
		t.Errorf("`pwd` with \"/\" as working directory: %v", err)
	}

	expected := "/"

	if expected != actual {
		t.Errorf("expected output of tool.New(\"pwd\").WithDir(\".\") to be %v: %v", expected, actual)
	}
}
