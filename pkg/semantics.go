package lexicalanalyser

import (
	"fmt"
)


type IDType string

const(
	Variable IDType = "variable"
	Function = "function"
)

type Id struct {
	typ IDType
	additionalInfo interface{}
}

type Scope struct {
	vars map[string]Id
	up *Scope
}

type Func struct {
	name string
	paramsVars *Scope
	typ VariableType
	bodyValid bool
}

type Var struct {
	name string
	typ VariableType
	value interface{}
}

type VariableType string

const(
	String = "string"
	Integer = "integer"
	Boolean = "boolean"
)

func stringToVariableType(typ string) VariableType {
	switch typ {
	case "string":
		return String
	case "int":
		return Integer
	case "bool":
		return Boolean
	default:
		return "unknown"
	}
}

func stringToBaseValue(typ string) interface{} {
	switch typ {
	case "string":
		return ""
	case "integer":
		return 0
	case "boolean":
		return false
	case "arr_string":
		return make([]string, 0)
	case "arr_integer":
		return make([]int, 0)
	case "arr_boolean":
		return make([]bool, 0)
	default:
		return "unknown"
	}
}

func (node *node32) checkSemantics(buffer string) error {
	//globalVars := make(map[string]Id) // key = variable name, value = variable type (Id)
	globalScope, tmpNode, err := node.up.getParamsVars(buffer)
	if err != nil {
		return fmt.Errorf("semantics error: %s", err)
	}
	fmt.Println(globalScope.vars)

	if rul3s[tmpNode.pegRule] == rul3s[ruleFUNCTIONS] {
		functionNode := tmpNode.up
		for functionNode != nil {
			//funcType := buffer[functionNode.up.next.next.begin:functionNode.up.next.next.end]
			//if rul3s[functionNode.up.next.pegRule] != rul3s[rulePARAMS_VARS] {
			//	funcType = buffer[functionNode.up.next.begin:functionNode.up.next.end]
			//}
			//
			//tmpId := Id{typ: Function, additionalInfo: Func{
			//	name: buffer[functionNode.up.begin: functionNode.up.end],
			//	paramsVars: functionScope,
			//	typ: stringToVariableType(funcType),
			//	bodyValid: true,
			//}}
			//globalScope.vars[buffer[functionNode.up.begin:functionNode.up.end]] = tmpId
			id, err := functionNode.validateFunction(buffer, globalScope)
			globalScope.vars[buffer[functionNode.up.begin:functionNode.up.end]] = id
			if err != nil {
				return fmt.Errorf("semantics error: %s", err)
			}
			functionNode = functionNode.next
		}
		tmpNode = tmpNode.next
	}
	//validate main body
	err = tmpNode.validateBody(buffer, globalScope)
	if err != nil {
		return fmt.Errorf("semantics error: main function invalid: %s", err)
	}
	return nil
}

func (node *node32) validateFunction(buffer string, scope *Scope) (Id, error) {
	functionNode := node
	_, isInScope := functionNode.isVarInScope(scope, buffer)
	if isInScope {
		return Id{}, fmt.Errorf("function name must be unique: %s", buffer[node.up.begin:node.up.end])
	}
	functionScope, _, err := functionNode.up.next.getParamsVars(buffer)
	functionScope.up = scope
	if err != nil {
		return Id{}, fmt.Errorf("function params vars: %s", err)
	}
	body := functionNode.up.next.next.next
	if rul3s[functionNode.up.next.next.pegRule] == rul3s[ruleBODY] {
		body = functionNode.up.next.next
	}
	if err := body.validateBody(buffer, functionScope); err != nil {
		return Id{}, fmt.Errorf("%s", err)
	}
	//todo: validate return in function, whether return returns predefined type
	//todo: return in if clause (and other body block) -> check whether if is in function
	//todo: functions must have unique names
	funcType := buffer[functionNode.up.next.next.begin:functionNode.up.next.next.end]
	if rul3s[functionNode.up.next.pegRule] != rul3s[rulePARAMS_VARS] {
		funcType = buffer[functionNode.up.next.begin:functionNode.up.next.end]
	}
	return Id{typ: Function, additionalInfo: Func{
		name: buffer[functionNode.up.begin: functionNode.up.end],
		paramsVars: functionScope,
		typ: stringToVariableType(funcType),
		bodyValid: true,
	}}, nil
}

func (node *node32) getParamsVars(buffer string) (*Scope, *node32, error) {
	vars := make(map[string]Id)
	tmpNode := node
	rule := rul3s[tmpNode.pegRule]
	for rule == rul3s[rulePARAMS_VARS] {
		typ := buffer[tmpNode.up.begin: tmpNode.up.end]
		varIDs := tmpNode.up.up.getVarIDs(buffer)
		for _, varID := range varIDs {
			if _, ok := vars[varID]; ok {
				return nil, node, fmt.Errorf("cannot define same variable ID more times: %s", varID)
			} else {
				vars[varID] = Id{typ: Variable,
					additionalInfo: Var{name: varID, typ: stringToVariableType(typ), value: stringToBaseValue(typ)},
				}
			}
		}
		tmpNode = tmpNode.next
		rule = rul3s[tmpNode.pegRule]
	}
	return &Scope{vars: vars, up: nil}, tmpNode, nil
}

func (node *node32) getVarIDs(buffer string) []string {
	tmpNode := node
	var res []string
	if rul3s[tmpNode.pegRule] == rul3s[ruleID] {
		for tmpNode != nil  {
			res = append(res, buffer[tmpNode.begin:tmpNode.end])
			tmpNode = tmpNode.next
		}
	}
	return res
}

func (node *node32) validateBody(buffer string, scope *Scope) error {
	if node == nil {
		return nil
	}

	if rul3s[node.pegRule] == rul3s[ruleBODY] && node.up != nil {
		bodyScope, statements, _ := node.up.getParamsVars(buffer)
		statement := statements.up
		for statement != nil {
			statementRule := rul3s[statement.up.pegRule]
			switch statementRule {
			case rul3s[ruleIF_STATEMENT]:
				bodyScope.up = scope
				if err := statement.up.validateIfStatement(buffer, bodyScope); err != nil {
					return fmt.Errorf("%s", err)
				}
			case rul3s[ruleWHILE_STATEMENT]:
				if err := statement.up.validateWhileStatement(buffer, scope); err != nil {
					return fmt.Errorf("%s", err)
				}
			case rul3s[ruleASSIGNMENT]:
				if err := statement.up.validateAssignment(buffer, scope); err != nil {
					return fmt.Errorf("%s", err)
				}
			//case rul3s[ruleFUNC_CALL]: // lebo viem modifikovat globalne premenne
			//	return validateFuncCall()
			//case rul3s[rulePRINT_STATEMENT]:
			//	return validatePrintStatement()
			}
			statement = statement.next
		}
	}
	return nil
}

//todo: vyriesit problem s rovnakym scope vo function params a function body
func (node *node32) validateIfStatement(buffer string, scope *Scope) error {
	rule := rul3s[node.pegRule]
	if rule == rul3s[ruleIF_STATEMENT] {
		bodyNode, err := node.up.checkBoolExpression(buffer, scope)
		if err != nil {
			return fmt.Errorf("if statement: %s", err)
		}
		err = bodyNode.validateBody(buffer, scope)
		if err != nil {
			return fmt.Errorf("if statement: %s", err)
		}
		//checking else clause
		if bodyNode.next != nil {
			err = bodyNode.next.up.validateBody(buffer, scope)
			if err != nil {
				return fmt.Errorf("if statement: %s", err)
			}
		}
	}
	return nil
}

func (node *node32) validateAssignment(buffer string, scope *Scope) error {
	if node == nil {
		return nil
	}

	if rul3s[node.pegRule] == rul3s[ruleASSIGNMENT] {
		varType, err := node.up.getValueType(buffer, scope)
		if err != nil {
			return fmt.Errorf("cannot determine variable type: %s", err)
		}

		//need to: assigning to array, checking type of indexing value
		//nedd to: check unindexed array assignment on left/right

		value := node.up.next
		valueType, err := value.up.getValueType(buffer, scope)
		if err != nil {
			return fmt.Errorf("could not get type of value: %s", err)
		}
		if varType != valueType {
			return fmt.Errorf("cannod assign %s value to %s variable: %s", valueType , varType , err)
		}

	}
	return nil
}

func (node *node32) validateWhileStatement(buffer string, scope *Scope) error {
	rule := rul3s[node.pegRule]
	if rule == rul3s[ruleWHILE_STATEMENT] {
		bodyNode, err := node.up.checkBoolExpression(buffer, scope)
		if err != nil {
			return fmt.Errorf("while statement: %s", err)
		}

		err = bodyNode.validateBody(buffer, scope)
		if err != nil {
			return fmt.Errorf("while statement: %s", err)
		}
	}
	return nil
}

func (node *node32) checkBoolExpression(buffer string, scope *Scope) (*node32, error) {
	if node == nil {
		return node, nil
	}
	rule := rul3s[node.pegRule]
	if rule == rul3s[ruleBOOL_EXPR_VALUE] {
		//check same types on every operation, whether it returns bool
		operandType, err := node.up.getValueType(buffer, scope)
		if err != nil {
			return node, err
		}
		tmp := node.next
		for rul3s[tmp.pegRule] != rul3s[ruleBODY] {
			if rul3s[tmp.pegRule] == rul3s[ruleBOOL_EXPR_VALUE] {
				valueType, _ := tmp.up.getValueType(buffer, scope)
				if operandType != valueType {
					return node, fmt.Errorf("cannot compare variables of different types: %s, %s",operandType, valueType)
				}
			}

			tmp = tmp.next
			rule = rul3s[tmp.pegRule]
		}

		return tmp, nil

	}
	return node, nil
}

func (node *node32) getValueType(buffer string, scope *Scope) (VariableType, error) {
	fmt.Println(buffer[node.begin:node.end])
	fmt.Println(node.pegRule)
	switch rul3s[node.pegRule] {
	case rul3s[ruleID]:
		//res := ""
		//if node.up != nil {
		//	//it's array item
		//	res = "arr_"
		//}
		tmpScope := scope
		for tmpScope != nil {
			typ, isInScope := node.isVarInScope(tmpScope, buffer)
			if isInScope {
				return typ, nil
			}
			tmpScope = tmpScope.up
		}

		return "", fmt.Errorf("variable was not declared before used: %s", buffer[node.begin:node.end])
		//find variable and return its type
	case rul3s[ruleTEXT]:
		return String, nil
	case rul3s[ruleINTEGER]:
		return Integer, nil
	case rul3s[ruleINT]:
		return Integer, nil
	case rul3s[ruleBOOLEAN]:
		return Boolean, nil
	case rul3s[ruleEXPRESSION]:
		exprType, err := node.up.validateExpression(buffer, scope)
		if err != nil {
			return "", err
		}
		return exprType, nil
	case rul3s[ruleFUNC_CALL]:
		//todo: check if function defined + get its type
		funcType, err := node.up.getValueType(buffer, scope)
		if err != nil {
			return "", err
		}
		return funcType, nil
	default:
		return "", fmt.Errorf("variable was not declared: %s", buffer[node.begin:node.end])
	}
}

func (node *node32) isVarInScope(scope *Scope, buffer string) (VariableType, bool) {
	variable := buffer[node.begin:node.end]
	if id, ok := scope.vars[variable]; ok {
		vari, ok := id.additionalInfo.(Var)
		if ok {
			return vari.typ, true
		}

		funct, ok := id.additionalInfo.(Func)
		if ok {
			return funct.typ, true
		}
	}
	return stringToVariableType("unknown"), false
}

func (node *node32) validateExpression(buffer string, scope *Scope) (VariableType, error) {
	if node == nil {
		return "", nil
	}

	if rul3s[node.pegRule] == rul3s[ruleEXPR_VALUE] {
		typeOfExprValue, err := node.up.getValueType(buffer, scope)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		}
		tmpNode := node.next
		for tmpNode != nil {
			if rul3s[tmpNode.pegRule] == rul3s[ruleOP] {
				op := buffer[tmpNode.begin:tmpNode.end]
				if typeOfExprValue == Boolean { // ID/func_call can be of type bool
					return "", fmt.Errorf("operation %s is not defined on boolean value", op)
				} else if typeOfExprValue == String && (op == "-" || op == "*" || op == "/") {
					return "", fmt.Errorf("operation %s is not defined on string value", op)
				}
			} else if rul3s[tmpNode.pegRule] == rul3s[ruleEXPR_VALUE] {
				nextType, err := tmpNode.up.getValueType(buffer, scope)
				if err != nil {
					return "", fmt.Errorf("%s", err)
				}
				if nextType != typeOfExprValue {
					return "", fmt.Errorf("cannot operate on values of different types: %s, %s", typeOfExprValue, nextType)
				}
			}
			tmpNode = tmpNode.next
		}
		//FUNC_CALL / TEXT / INTEGER / ID -> (ID/ func_call moze mat typ bool)
		//+: b+b
		//*: s*s, b*b
		//-: s-s, b-b
		///: s/s, b/b, (!) i/i- > check division by 0 when generating + vysledny typ
		return typeOfExprValue, nil
	}

	return "", fmt.Errorf("validateExpression called on wrong node: %v", node.pegRule)
}