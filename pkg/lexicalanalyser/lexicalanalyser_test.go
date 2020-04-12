package lexicalanalyser

import (
	"fmt"
	"strings"
	"testing"
)

type mockLexA struct {
	Program []string
	Current int
}

type mockLexicalAnalyser interface {
	LexicalAnalyser
	GetProgram() []string
}

func newMockLexicalAnalyzer(program string) mockLexicalAnalyser {
	return &mockLexA{Program: strings.Fields(program), Current: 0}
}

func (l mockLexA) GetToken() (string, error) {
	if l.Current == len(l.Program)-1 {
		return " ", fmt.Errorf("no other string available")
	}
	res := l.Program[l.Current]

	l.Current++

	return res, nil
}

func (l mockLexA) NextToken() (string, error) {
	if l.Current == len(l.Program)-1 {
		return " ", fmt.Errorf("no other rune available")
	}
	return l.Program[l.Current], nil
}

func (l mockLexA) GetProgram() []string {
	return l.Program
}

func TestGetToken(t *testing.T) {
	t.Run("test getToken with only spaces", func(t *testing.T) {
		la := newMockLexicalAnalyzer("globals endglobals main { ; return void } endmain")
		assertProgramLength(t, len(la.GetProgram()), 9)
	})

	t.Run("test getToken with spaces, tabs, newlines", func(t *testing.T) {
		la := newMockLexicalAnalyzer(`globals 	

			
	endglobals main { ; return voidV } endmain`)
		assertProgramLength(t, len(la.GetProgram()), 9)
	})
}

func assertProgramLength(t *testing.T, actual, expected int) {
	t.Helper()
	if actual != expected {
		t.Errorf("expected %v actual %v", expected, actual)
	}
}
