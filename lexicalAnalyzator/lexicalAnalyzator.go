package lexicalAnalyzator

import (
	"fmt"
	"strings"
)

type LexicalAnalyzer interface {
	GetToken() (string, error)
	NextToken() (string, error)
}

type lexA struct {
	program []string
	current int
}

func NewLexicalAnalyzer(program string) LexicalAnalyzer {
	return &lexA{program: strings.Fields(program), current: 0}
}

func (l *lexA) NextToken() (string, error) {
	if l.current == len(l.program) {
		return " ", fmt.Errorf("no other string available")
	}
	res := l.program[l.current]
	l.current++

	return res, nil
}

func (l *lexA) GetToken() (string, error) {
	if l.current == len(l.program) {
		return " ", fmt.Errorf("no other rune available")
	}

	return l.program[l.current], nil
}