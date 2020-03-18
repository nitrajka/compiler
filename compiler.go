package main

type Program struct {

}

type Compiler interface {
	Compile() Program
}

type compiler struct {
	la LexicalAnalyzer
	stack []rune
}

func NewCompiler(la LexicalAnalyzer) Compiler {
	return &compiler{la:la}
}

func (c compiler) Compile() Program {
	token, err := c.la.GetToken()
	for err == nil {
		token, err = c.la.GetToken()
	}
	return Program{}
}

func (c compiler) CompileStatement() {

}