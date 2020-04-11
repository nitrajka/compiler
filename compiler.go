package main

import (
	"fmt"
	"github.com/nitrajka/compiler/lexicalAnalyzator"
)

type Program struct {

}

type Compiler interface {
	Compile() (Program, error)
}

type compiler struct {
	la    lexicalAnalyzator.LexicalAnalyzer
	stack []string
}

func NewCompiler(la lexicalAnalyzator.LexicalAnalyzer) Compiler {
	return &compiler{la: la, stack: []string{}}
}

func (c compiler) Compile() (Program, error) {
	token, err := c.la.NextToken()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v\n", err)
	}
	if token != "globals" {
		return Program{}, fmt.Errorf("compile error, ROOT: unexpected token %v\n", token)
	}

	err = c.compileParamVars()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v\n", err)
	}

	token, err = c.la.NextToken()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v\n", err)
	}
	if token != "endglobals" {
		return Program{}, fmt.Errorf("compile error, ROOT: unexpected token %v\n", token)
	}

	token, err = c.la.NextToken()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v\n", err)
	}
	if token == "func" {
		err := c.compileFunctions()
		if err != nil {
			return Program{}, fmt.Errorf("compile error, ROOT: %v\n", err)
		}
	} else if token == "main" {
		err := c.compileBody()
		if err != nil {
			return Program{}, fmt.Errorf("compile error, ROOT: %v\n", err)
		}
	} else {
		return Program{}, fmt.Errorf("compile error, ROOT: unexpected token: %v\n", token)
	}

	token, err = c.la.NextToken()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v\n", err)
	}
	if token != "endmain" {
		return Program{}, fmt.Errorf("compile error, ROOT: unexpected token %v\n", token)
	}

	return Program{}, nil
}

func (c compiler) compileParamVars() error {
	return nil
}

func (c compiler) compileFunctions() error {
	return nil
}

func (c compiler) compileBody() error {
	token, err := c.la.NextToken()
	if err != nil {
		return fmt.Errorf("compile error, BODY: %v\n", err)
	}
	if token != "{" {
		return fmt.Errorf("compile error, BODY: unexpected token %v\n", token)
	}

	err = c.compileParamVars()
	if err != nil {
		return fmt.Errorf("compile error, BODY: %v\n", err)
	}

	token, err = c.la.NextToken()
	if err != nil {
		return fmt.Errorf("compile error, BODY: %v\n", err)
	}
	if token != ";" {
		return fmt.Errorf("compile error, BODY: unexpected token %v\n", token)
	}

	c.compileReturnClause()

	token, err = c.la.NextToken()
	if err != nil {
		return fmt.Errorf("compile error, BODY: %v\n", err)
	}
	if token != "}" {
		return fmt.Errorf("compile error, BODY: unexpected token %v\n", token)
	}

	return nil
}

func (c compiler) compileReturnClause() error {
	token, err := c.la.NextToken()
	if err != nil {
		return fmt.Errorf("compile error, RETURN: %v\n", err)
	}
	if token != "return" {
		return fmt.Errorf("compile error, RETURN: %v\n", err)
	}

	return c.compileValue()
}

func (c compiler) compileValue() error {
	token, err := c.la.NextToken()
	if err != nil {
		return fmt.Errorf("compile error, RETURN: %v\n", err)
	}

	switch token {
	case "voidV":
		return nil
	default:
		return fmt.Errorf("compile error, VALUE: unexpected token %v\n", token)
	}

	return nil
}