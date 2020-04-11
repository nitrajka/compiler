package main

import (
	"github.com/nitrajka/compiler/lexicalAnalyzator"
	"testing"
)

func TestCompile(t *testing.T) {
	t.Run("simplest correct program", func(t *testing.T) {
		la := lexicalAnalyzator.NewLexicalAnalyzer("globals endglobals main { ; return void } endmain")
		compiler := NewCompiler(la)
		_, err := compiler.Compile()
		assertError(t, err, nil)
	})
}

func assertError(t *testing.T, actual, expected error) {
	t.Helper()
	if actual != expected {
		t.Errorf("expected %v actual %v", expected, actual)
	}
}
