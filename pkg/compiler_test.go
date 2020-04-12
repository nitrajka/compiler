package pkg

import (
	"testing"

	"github.com/nitrajka/compiler/pkg/lexicalanalyser"
)

func TestCompile(t *testing.T) {
	t.Run("simplest correct program", func(t *testing.T) {
		la := lexicalanalyser.NewLexicalAnalyzer("globals endglobals main { ; return void } endmain")
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
