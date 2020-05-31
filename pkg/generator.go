package lexicalanalyser

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"os"
	"strconv"
)

func (node *node32) Generate(buffer string, to string) error {
	//node.PrettyPrint(os.Stdout, buffer)

	//todo: unused variables

	f := jen.NewFile("main")

	tmpNode, paramsVars := node.up.generateParamsVars(buffer)
	for _, paramVar := range paramsVars {
		f.Var().Add(paramVar.name, paramVar.typ)
	}

	//tmpNode = tmpNode.generateFunctions(buffer, f)

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

type vr struct {
	name *jen.Statement
	typ *jen.Statement
}

func (node *node32) generateParamsVars(buffer string) (*node32, []vr) {
	tmpNode := node
	var res []vr
	for tmpNode.pegRule == rulePARAMS_VARS {

		tmpId := tmpNode.up.up
		for tmpId != nil {
			res = append(res, vr{name: jen.Id(buffer[tmpId.begin:tmpId.end]), typ: defineVariable(buffer[tmpNode.up.begin:tmpNode.up.end]) })
			tmpId = tmpId.next
		}

		tmpNode = tmpNode.next
	}
	return tmpNode, res
}

func defineVariable(typ string) *jen.Statement {
	switch stringToVariableType(typ) {
	case String:
		return jen.String()
	case Integer:
		return jen.Int()
	case Boolean:
		return jen.Bool()
	}
	return nil
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
		statementsAst, p := node.up.generateParamsVars(buffer)
		fmt.Println(p)
		statement := statementsAst.up
		var statements []*jen.Statement

		var value *node32
		for statement != nil {
			if statement.up.pegRule == rulePRINT_STATEMENT {
				value = statement.up.up.up
				generatedValue := value.generateOperand(buffer)
				statements = append(statements, jen.Qual("fmt", "Println").Call(generatedValue))
			} else if statement.up.pegRule == ruleIF_STATEMENT {

			} else if statement.up.pegRule == ruleWHILE_STATEMENT {
				boolExpr, body := statement.up.up.getBoolExprValue(buffer)
				body.generateBody(buffer, jen.For(boolExpr...))
			} else if statement.up.pegRule == ruleASSIGNMENT {
				value := statement.up.up.next


				if value.up.pegRule == ruleEXPRESSION {
					generatedValue := value.up.getExprValue(buffer)
					statements = append(statements, statement.up.up.generateOperand(buffer).Op("=").Add(generatedValue...))
				} else {
					generatedValue := value.up.generateOperand(buffer)
					statements = append(statements, statement.up.up.generateOperand(buffer).Op("=").Add(generatedValue))
				}

			}
			statement = statement.next
		}

		//var statementsToBlock []*jen.Statement
		//tmpNode := node.generateParamsVars()


		var code []jen.Code
		for _, paramVar := range p {
			code = append(code, jen.Var().Add(paramVar.name, paramVar.typ))
		}
		for _, s := range statements {
			code = append(code, s)
		}
		s.Block(code...)
		//f.Func().Id(funcName).Params().Block(
		//	jen.Qual("fmt", "Println").Call(jen.Lit("Hello, world")),
		//)
	}

}

func (node *node32) getExprValue(buffer string) []jen.Code {
	var res []jen.Code

	if node.pegRule == ruleEXPRESSION {
		var leftOp *jen.Statement
		var op string
		tmpNode := node.up
		for tmpNode != nil && (tmpNode.pegRule == ruleEXPR_VALUE || tmpNode.pegRule == ruleOP) {
			if tmpNode.pegRule == ruleEXPR_VALUE {
				if leftOp == nil {
					leftOp = tmpNode.up.generateOperand(buffer)
				} else {
					right := tmpNode.up.generateOperand(buffer)
					res = append(res, leftOp.Op(op).Add(right) )
					leftOp = right
				}
			} else if tmpNode.pegRule == ruleOP {
				op = buffer[tmpNode.begin:tmpNode.end]
			}
			tmpNode = tmpNode.next
		}
	}

	return res
}

func (node *node32) getBoolExprValue(buffer string) ([]jen.Code, *node32) {
	var res []jen.Code
	tmpNode := node

	if node.pegRule == ruleBOOL_EXPR_VALUE {
		var leftOp *jen.Statement
		var op string
		for tmpNode != nil && (tmpNode.pegRule == ruleBOOL_EXPR_VALUE || tmpNode.pegRule == ruleBOOL_OP) {
			if tmpNode.pegRule == ruleBOOL_EXPR_VALUE {
				if leftOp != nil { // standing on right operand -> generate operation

					if tmpNode.up.pegRule == ruleID {
						res = append(res, leftOp.Op(op).Id(buffer[tmpNode.up.begin:tmpNode.up.end]))
					} else if tmpNode.up.pegRule == ruleBOOLEAN {
						if buffer[tmpNode.up.begin:tmpNode.up.end] == "true" {
							res = append(res, leftOp.Op(op).True())
						} else {
							res = append(res, leftOp.Op(op).False())
						}
					} else if tmpNode.up.pegRule == ruleINTEGER || tmpNode.up.pegRule == ruleTEXT {
						res = append(res, leftOp.Op(op).Lit(buffer[tmpNode.up.begin:tmpNode.up.end]))
					}
					leftOp = tmpNode.up.generateOperand(buffer)
				} else {
					leftOp = tmpNode.up.generateOperand(buffer)
				}
			} else if tmpNode.pegRule == ruleBOOL_OP {
				op = buffer[tmpNode.begin:tmpNode.end]
			}
			tmpNode = tmpNode.next
		}
	}
	return res, tmpNode
}

func (node *node32) generateOperand(buffer string) *jen.Statement {
	if node.pegRule == ruleID {
		return jen.Id(buffer[node.begin:node.end])
	} else if node.pegRule == ruleBOOLEAN {
		return node.generateBool(buffer[node.begin:node.end])
	} else if node.pegRule == ruleINTEGER {
		num, _ := strconv.Atoi(buffer[node.begin:node.end])
		return jen.Lit(num)
	} else if node.pegRule == ruleTEXT {
		return jen.Lit(buffer[node.begin:node.end])
	}
	return nil
}

func (node *node32) generateBool(bl string) *jen.Statement {
	if bl == "true" {
		return jen.True()
	}
	return jen.False()
}