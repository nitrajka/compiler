package pkg

import (
	"fmt"

	"github.com/nitrajka/compiler/pkg/lexicalanalyser"
)

type Program struct {
}

type Compiler struct {
	la    lexicalanalyser.LexicalAnalyser
	stack []string
}

func NewCompiler(la lexicalanalyser.LexicalAnalyser) *Compiler {
	return &Compiler{la: la, stack: []string{}}
}

func (c Compiler) Compile() (Program, error) {
	token, err := c.la.NextToken()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v", err)
	}
	if token != "globals" {
		return Program{}, fmt.Errorf("compile error, ROOT: unexpected token %v", token)
	}

	err = c.compileParamVars()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v", err)
	}

	token, err = c.la.NextToken()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v", err)
	}
	if token != "endglobals" {
		return Program{}, fmt.Errorf("compile error, ROOT: unexpected token %v", token)
	}

	token, err = c.la.NextToken()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v", err)
	}
	if token == "func" {
		err = c.compileFunctions()
		if err != nil {
			return Program{}, fmt.Errorf("compile error, ROOT: %v", err)
		}
	} else if token == "main" {
		err = c.compileBody()
		if err != nil {
			return Program{}, fmt.Errorf("compile error, ROOT: %v", err)
		}
	} else {
		return Program{}, fmt.Errorf("compile error, ROOT: unexpected token: %v", token)
	}

	token, err = c.la.NextToken()
	if err != nil {
		return Program{}, fmt.Errorf("compile error, ROOT: %v", err)
	}
	if token != "endmain" {
		return Program{}, fmt.Errorf("compile error, ROOT: unexpected token %v", token)
	}

	return Program{}, nil
}

func (c Compiler) compileParamVars() error {
	return nil
}

func (c Compiler) compileFunctions() error {
	return nil
}

func (c Compiler) compileBody() error {
	if token, err := c.la.NextToken(); err != nil {
		return fmt.Errorf("compile error, BODY: %v", err)
	} else if token != "{" {
		return fmt.Errorf("compile error, BODY: unexpected token %v", token)
	}

	if err := c.compileParamVars(); err != nil {
		return fmt.Errorf("compile error, BODY: %v", err)
	}

	if token, err := c.la.NextToken(); err != nil {
		return fmt.Errorf("compile error, BODY: %v", err)
	} else if token != ";" {
		return fmt.Errorf("compile error, BODY: unexpected token %v", token)
	}

	if err := c.compileReturnClause(); err != nil {
		return fmt.Errorf("todo%v", err)
	}

	if token, err := c.la.NextToken(); err != nil {
		return fmt.Errorf("compile error, BODY: %v", err)
	} else if token != "}" {
		return fmt.Errorf("compile error, BODY: unexpected token %v", token)
	}

	return nil
}

func (c Compiler) compileReturnClause() error {
	token, err := c.la.NextToken()
	if err != nil {
		return fmt.Errorf("compile error, RETURN: %v", err)
	}
	if token != "return" {
		return fmt.Errorf("compile error, RETURN: %v", err)
	}

	return c.compileValue()
}

func (c Compiler) compileValue() error {
	token, err := c.la.NextToken()
	if err != nil {
		return fmt.Errorf("compile error, RETURN: %v", err)
	}

	switch token {
	case "void":
		return nil
	default:
		return fmt.Errorf("compile error, VALUE: unexpected token %v", token)
	}
}