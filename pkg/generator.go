package lexicalanalyser

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"os"
)

func (node *node32) Generate(buffer string, to string) error {
	//node.PrettyPrint(os.Stdout, buffer)

	f := jen.NewFile("main")

	tmpNode := node.generateParamsVars(buffer, f)

	tmpNode = tmpNode.generateFunctions(buffer, f)

	tmpNode.generateBody(buffer, f.Func().Id("main").Params())

	file, err := os.Create(to)
	if err != nil {
		return err
	}
	nb, err2 := fmt.Fprintf(file, "%#v", f)
	if err2 != nil {
		return err2
	}
	_ = nb
	//fmt.Printf("%#v\n", f)
	//fmt.Printf("%d bytes written\n", nb)
	return nil
}

func (node *node32) generateParamsVars(buffer string, f *jen.File) *node32 {
	tmpNode := node.up
	for tmpNode.pegRule == rulePARAMS_VARS {

		tmpId := tmpNode.up.up
		for tmpId != nil {
			defineVariable(buffer[tmpNode.up.begin:tmpNode.up.end], buffer[tmpId.begin:tmpId.end], f)
			tmpId = tmpId.next
		}

		tmpNode = tmpNode.next
	}
	return tmpNode
}

func defineVariable(typ string, name string, f *jen.File) {
	switch stringToVariableType(typ) {
	case String:
		f.Var().Id(name).String()
	case Integer:
		f.Var().Id(name).Int()
	case Boolean:
		f.Var().Id(name).Bool()
	}
}

func (node *node32) generateFunctions(buffer string, f *jen.File) *node32 {
	tmpNode := node
	if tmpNode.pegRule == ruleFUNCTIONS {

		tmpFunction := tmpNode.up
		for tmpFunction != nil {
			//todo: generate function
			//tmpStatements := tmpFunction.up.generateParamsVars()
			//f.Func().Id(buffer[tmpFunction.up.begin:tmpFunction.up.end]).Block()
			tmpFunction.generateBody(buffer, f.Func().Id("a").Params())
			tmpFunction = tmpFunction.next
		}
		tmpNode = tmpNode.next
	}
	return tmpNode
}

func (node *node32) generateBody(buffer string, s *jen.Statement) {
	if node.pegRule == ruleBODY {
		statementsAst := node.up
		statement := statementsAst.up
		var value *node32
		if statement.up.pegRule == rulePRINT_STATEMENT {
			value = statement.up.up.up
		}

		//var statementsToBlock []*jen.Statement
		//tmpNode := node.generateParamsVars()
		var statements []*jen.Statement

		statements = append(statements, jen.Qual("fmt", "Println").Call(jen.Lit(buffer[value.begin:value.end])))

		var code []jen.Code
		for _, s := range statements {
			code = append(code, s)
		}
		s.Block(code...)

		//f.Func().Id(funcName).Params().Block(
		//	jen.Qual("fmt", "Println").Call(jen.Lit("Hello, world")),
		//)
	}
}