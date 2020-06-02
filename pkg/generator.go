package lexicalanalyser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dave/jennifer/jen"
)

func (node *node32) Generate(buffer, to string) error {
	// node.PrettyPrint(os.Stdout, buffer)

	f := jen.NewFile("main")

	tmpNode, paramsVars := node.up.generateParamsVars(buffer)
	for _, paramVar := range paramsVars {
		f.Var().Add(paramVar.name, paramVar.typ)
	}

	tmpNode = tmpNode.generateFunctions(buffer, f)

	tmpNode.generateBody(buffer, f.Func().Id("main").Params(), paramsVars)

	file, err := os.Create(to)
	if err != nil {
		return err
	}

	if _, err2 := fmt.Fprintf(file, "%#v", f); err2 != nil {
		return err2
	}

	return nil
}

type vr struct {
	name *jen.Statement
	typ  *jen.Statement
}

func (node *node32) generateParamsVars(buffer string) (*node32, []vr) {
	if node == nil {
		return node, make([]vr, 0)
	}

	var res []vr

	tmpNode := node
	for tmpNode.pegRule == rulePARAMS_VARS {
		tmpID := tmpNode.up.up
		for tmpID != nil {
			res = append(res, vr{
				name: jen.Id(buffer[tmpID.begin:tmpID.end]),
				typ:  defineVariable(buffer[tmpNode.up.begin:tmpNode.up.end]),
			})
			tmpID = tmpID.next
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
			funcName := buffer[tmpFunction.up.begin:tmpFunction.up.end]
			typ, params := tmpFunction.up.next.generateParamsVars(buffer)

			var tmpParams []jen.Code
			for _, p := range params {
				tmpParams = append(tmpParams, p.name.Add(p.typ))
			}

			typeStmnt := defineVariable(buffer[typ.begin:typ.end])
			fu := jen.Func().Id(funcName).Params(tmpParams...).Add(typeStmnt)
			typ.next.generateBody(buffer, fu, nil)
			f.Add(fu)

			tmpFunction = tmpFunction.next
		}

		tmpNode = tmpNode.next
	}

	return tmpNode
}

func (node *node32) generateBody(buffer string, s *jen.Statement, globalVars []vr) {
	if node.pegRule != ruleBODY {
		return
	}

	statementsAst, p := node.up.generateParamsVars(buffer)

	var statements []*jen.Statement
	var value *node32

	if statementsAst != nil && statementsAst.pegRule == ruleSTATEMENTS {
		statement := statementsAst.up
		for statement != nil {
			switch statement.up.pegRule {
			case rulePRINT_STATEMENT:
				if statement.up.up != nil {
					value = statement.up.up
					generatedValue := value.up.generateOperand(buffer)
					statements = append(statements, jen.Qual("fmt", "Println").Call(generatedValue))
				} else { // empty print, printing newline
					statements = append(statements, jen.Qual("fmt", "Println").Call())
				}
			case ruleIF_STATEMENT:
				boolExpr, body := statement.up.up.getBoolExprValue(buffer)
				k := jen.If(boolExpr...)

				body.generateBody(buffer, k, nil)

				if body.next != nil { // elseclause exists
					body.next.up.generateBody(buffer, k.Else(), nil)
				}

				statements = append(statements, k)
			case ruleWHILE_STATEMENT:
				boolExpr, body := statement.up.up.getBoolExprValue(buffer)
				k := jen.For(boolExpr...)

				body.generateBody(buffer, k, nil)
				statements = append(statements, k)
			case ruleASSIGNMENT:
				value := statement.up.up.next
				leftSideOfAssignment := statement.up.up.generateOperand(buffer).Op("=")
				switch value.up.pegRule {
				case ruleEXPRESSION:
					generatedValue := value.up.getExprValue(buffer)
					statements = append(statements, leftSideOfAssignment.Add(generatedValue...))
				case ruleFUNC_CALL:
					var params []jen.Code

					tmpID := value.up.up.next
					for tmpID != nil {
						params = append(params, jen.Id(buffer[tmpID.begin:tmpID.end]))
						tmpID = tmpID.next
					}
					statements = append(statements, leftSideOfAssignment.Id(buffer[value.up.up.begin:value.up.up.end]).Call(params...))
				default:
					generatedValue := value.up.generateOperand(buffer)
					statements = append(statements, leftSideOfAssignment.Add(generatedValue))
				}
			case ruleFUNC_CALL:
				var params []jen.Code

				tmpID := statement.up.up.next
				for tmpID != nil {
					params = append(params, jen.Id(buffer[tmpID.begin:tmpID.end]))
					tmpID = tmpID.next
				}
				statements = append(statements, jen.Id(buffer[statement.up.up.begin:statement.up.up.end]).Call(params...))
			}

			statement = statement.next
		}

		statementsAst = statementsAst.next
	}

	if statementsAst != nil && statementsAst.pegRule == ruleRETURN_CLAUSE {
		value := statementsAst.up
		if value.up != nil { // is not void
			switch value.up.pegRule {
			case ruleEXPRESSION:
				code := value.up.getExprValue(buffer)
				statements = append(statements, jen.Return(code...))
			case ruleFUNC_CALL:
				var params []jen.Code
				tmpID := value.up.up.next
				for tmpID != nil {
					params = append(params, jen.Id(buffer[tmpID.begin:tmpID.end]))
					tmpID = tmpID.next
				}
				statements = append(statements, jen.Return().Id(buffer[value.up.up.begin:value.up.up.end]).Call(params...))
			default:
				stmnt := value.up.generateOperand(buffer)
				statements = append(statements, jen.Return(stmnt))
			}
		} else {
			statements = append(statements, jen.Return())
		}
	}

	var code []jen.Code
	for _, paramVar := range p {
		code = append(
			code,
			jen.Var().Add(paramVar.name, paramVar.typ),
			jen.Id("_").Op("=").Add(paramVar.name),
		)
	}

	// _ = varName -> to make sure there is not an unused variable
	for _, paramVar := range globalVars {
		code = append(code, jen.Id("_").Op("=").Add(paramVar.name))
	}

	for _, s := range statements {
		code = append(code, s)
	}

	s.Block(code...)
}

func (node *node32) getExprValue(buffer string) []jen.Code {
	var res []jen.Code

	if node.pegRule != ruleEXPRESSION {
		return res
	}

	var leftOp *jen.Statement
	var op string

	tmpNode := node.up
	for tmpNode != nil && (tmpNode.pegRule == ruleEXPR_VALUE || tmpNode.pegRule == ruleOP) {
		if tmpNode.pegRule == ruleEXPR_VALUE {
			if leftOp == nil {
				leftOp = tmpNode.up.generateOperand(buffer)
			} else {
				right := tmpNode.up.generateOperand(buffer)
				leftOp = leftOp.Add(right)
			}
		} else if tmpNode.pegRule == ruleOP {
			op = buffer[tmpNode.begin:tmpNode.end]
			leftOp = leftOp.Op(op)
		}

		tmpNode = tmpNode.next
	}

	res = append(res, leftOp)

	return res
}

func (node *node32) getBoolExprValue(buffer string) ([]jen.Code, *node32) {
	var res []jen.Code

	if node.pegRule != ruleBOOL_EXPR_VALUE {
		return res, node
	}

	var op string
	var wasOp bool
	var leftOp *jen.Statement
	tmpNode := node

	for tmpNode != nil && (tmpNode.pegRule == ruleBOOL_EXPR_VALUE || tmpNode.pegRule == ruleBOOL_OP) {
		if tmpNode.pegRule == ruleBOOL_EXPR_VALUE {
			if leftOp == nil { // standing on right operand -> generate operation
				leftOp = tmpNode.up.generateOperand(buffer)
			} else {
				var tmp *jen.Statement

				switch tmpNode.up.pegRule {
				case ruleID:
					leftOp = leftOp.Id(buffer[tmpNode.up.begin:tmpNode.up.end])
					tmp = jen.Id(buffer[tmpNode.up.begin:tmpNode.up.end])
				case ruleBOOLEAN:
					bl := tmpNode.up.generateBool(buffer[tmpNode.up.begin:tmpNode.up.end])
					leftOp = leftOp.Add(bl)
					tmp = bl
				case ruleINTEGER:
					num, _ := strconv.Atoi(buffer[tmpNode.up.begin:tmpNode.up.end])
					leftOp = leftOp.Lit(num)
					tmp = jen.Lit(num)
				case ruleTEXT:
					leftOp = leftOp.Lit(buffer[tmpNode.up.begin:tmpNode.up.end])
					tmp = jen.Lit(buffer[tmpNode.up.begin:tmpNode.up.end])
				}

				if wasOp && tmpNode.next != nil && tmpNode.next.pegRule == ruleBOOL_OP {
					leftOp = leftOp.Op("&&").Add(tmp)
				}
			}
		} else if tmpNode.pegRule == ruleBOOL_OP {
			op = buffer[tmpNode.begin:tmpNode.end]
			leftOp = leftOp.Op(op)
			wasOp = true
		}

		tmpNode = tmpNode.next
	}

	res = append(res, leftOp)

	return res, tmpNode
}

func (node *node32) generateOperand(buffer string) *jen.Statement {
	switch node.pegRule {
	case ruleID:
		return jen.Id(buffer[node.begin:node.end])
	case ruleBOOLEAN:
		return node.generateBool(buffer[node.begin:node.end])
	case ruleINTEGER:
		num, _ := strconv.Atoi(buffer[node.begin:node.end])
		return jen.Lit(num)
	case ruleTEXT:
		if buffer[node.begin:node.end] == "\"\"" {
			return jen.Lit("")
		}
		return jen.Lit(buffer[node.begin:node.end])
	case ruleFUNC_CALL:
		var params []jen.Code
		tmpID := node.up.next
		for tmpID != nil {
			params = append(params, jen.Id(buffer[tmpID.begin:tmpID.end]))
			tmpID = tmpID.next
		}
		return jen.Id(buffer[node.up.begin:node.up.end]).Call(params...)
	case ruleEXPRESSION:
		code := node.getExprValue(buffer)
		return jen.Add(code...)
	}

	return nil
}

func (node *node32) generateBool(bl string) *jen.Statement {
	if bl == "true" {
		return jen.True()
	}

	return jen.False()
}