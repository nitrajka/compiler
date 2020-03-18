package main

import "fmt"

type LexicalAnalyzer interface {
	GetToken() (rune, error)
	NextToken() (rune, error)
}

type lexA struct {
	program []rune
	current int
}

func NewLexicalAnalyzer(program []rune) LexicalAnalyzer {
	return &lexA{program:program, current:0}
}

func (l *lexA) GetToken() (rune, error) {
	if l.current == len(l.program) {
		return ' ', fmt.Errorf("no other rune available")
	}
	res := l.program[l.current]
	l.current++

	return res, nil
}

func (l *lexA) NextToken() (rune, error) {
	if l.current >= len(l.program)-1 {
		return ' ', fmt.Errorf("no other rune available")
	}
	return l.program[l.current+1], nil
}