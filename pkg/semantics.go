package lexicalanalyser

import (
	"fmt"
)

type VariableKind string

const (
	Variable VariableKind = "variable"
	Function VariableKind = "function"
)

type ID struct {
	variableKind   VariableKind
	variableType   VariableType
	name           string
	additionalInfo interface{}
}

type Scope struct {
	vars            map[string]ID
	up              *Scope
	currentFunction string
}

type Func struct {
	paramsVars  *Scope
	paramsOrder []VariableType
	bodyValid   bool
}

type Var struct {
	value interface{}
}

type VariableType string

const (
	String  = "string"
	Integer = "integer"
	Boolean = "boolean"
	Void    = "void"
	Unknown = "unknown"
)

func stringToVariableType(typ string) VariableType {
	switch typ {
	case "string": // nolint
		return String
	case "int":
		return Integer
	case "bool":
		return Boolean
	case "void": // nolint
		return Void
	default:
		return Unknown
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
	case "void":
		return struct{}{}
	default:
		return Unknown
	}
}

func (node *node32) checkSemantics(buffer string) error {
	globalScope, tmpNode, _, err := node.up.getParamsVars(buffer, nil)
	if err != nil {
		return fmt.Errorf("semantics error: %s", err)
	}

	if tmpNode.pegRule == ruleFUNCTIONS {
		functionNode := tmpNode.up

		// declare all functions into global scope
		for functionNode != nil {
			id, err2 := functionNode.declareFunction(buffer, globalScope)
			if err2 != nil {
				return fmt.Errorf("semantics error: %s", err2)
			}
			globalScope.vars[buffer[functionNode.up.begin:functionNode.up.end]] = id
			functionNode = functionNode.next
		}

		// define all previously declared functions
		functionNode = tmpNode.up
		for functionNode != nil {
			id, err2 := functionNode.validateFunction(buffer, globalScope)
			globalScope.vars[buffer[functionNode.up.begin:functionNode.up.end]] = id // additionalInfo changes
			if err2 != nil {
				return fmt.Errorf("semantics error: %s", err2)
			}
			functionNode = functionNode.next
		}
		tmpNode = tmpNode.next
	}

	// validate main body
	globalScope.currentFunction = "main"
	err = tmpNode.validateBody(buffer, globalScope, false)
	if err != nil {
		return fmt.Errorf("semantics error: main function invalid: %s", err)
	}

	return nil
}

func (node *node32) declareFunction(buffer string, scope *Scope) (ID, error) {
	_, isInScope := node.up.isVarInScope(scope, buffer)
	if isInScope {
		return ID{}, fmt.Errorf("function name must be unique: %s", buffer[node.up.begin:node.up.end])
	}

	functionScope, nodeType, paramsOrder, err := node.up.next.getParamsVars(buffer, nil)
	if err != nil {
		return ID{}, fmt.Errorf("function params vars: %s", err)
	}
	functionScope.up = scope
	functionScope.currentFunction = buffer[node.up.begin:node.up.end]

	funcType := buffer[nodeType.begin:nodeType.end]

	return ID{
		variableKind:   Function,
		variableType:   stringToVariableType(funcType),
		name:           buffer[node.up.begin:node.up.end],
		additionalInfo: Func{paramsVars: functionScope, paramsOrder: paramsOrder, bodyValid: false},
	}, nil
}

func (node *node32) validateFunction(buffer string, scope *Scope) (ID, error) {
	functionNode := node
	functionScope, _, paramsOrder, err := functionNode.up.next.getParamsVars(buffer, nil)
	if err != nil {
		return ID{}, fmt.Errorf("function params vars: %s", err)
	}
	functionScope.up = scope
	functionScope.currentFunction = buffer[functionNode.up.begin:functionNode.up.end]
	body := functionNode.up.next.next.next
	if functionNode.up.next.next.pegRule == ruleBODY {
		body = functionNode.up.next.next
	}
	if err := body.validateBody(buffer, functionScope, true); err != nil {
		return ID{}, fmt.Errorf("%s", err)
	}
	//todo: validate whether function with return type has return statement (only void does not have to have)
	//todo: 2 unknown vars in assignment (or else), check for unknown
	//todo: validate print statement
	//todo: generovanie kodu
	tmpNode := node.up.next
	for tmpNode.pegRule == rulePARAMS_VARS {
		tmpNode = tmpNode.next
	}
	funcType := buffer[tmpNode.begin:tmpNode.end]

	return ID{
		variableKind: Function,
		name:         buffer[functionNode.up.begin:functionNode.up.end],
		variableType: stringToVariableType(funcType),
		additionalInfo: Func{
			paramsVars:  functionScope,
			paramsOrder: paramsOrder,
			bodyValid:   true,
		}}, nil
}

func (node *node32) getParamsVars(buffer string, scope *Scope) (*Scope, *node32, []VariableType, error) {
	if scope == nil {
		scope = &Scope{vars: make(map[string]ID)}
	}

	var paramsOrder []VariableType

	tmpNode := node

	rule := tmpNode.pegRule
	for rule == rulePARAMS_VARS {
		typ := buffer[tmpNode.up.begin:tmpNode.up.end]
		varIDs := tmpNode.up.up.getVarIDs(buffer)

		for _, varID := range varIDs {
			paramsOrder = append(paramsOrder, stringToVariableType(typ))

			if _, ok := scope.vars[varID]; ok {
				return nil, node, nil, fmt.Errorf("cannot define same variable ID more times in the same scope: %s", varID)
			}
			scope.vars[varID] = ID{variableKind: Variable, name: varID, variableType: stringToVariableType(typ),
				additionalInfo: Var{value: stringToBaseValue(typ)},
			}
		}
		tmpNode = tmpNode.next
		rule = tmpNode.pegRule
	}

	return scope, tmpNode, paramsOrder, nil
}

func (node *node32) getVarIDs(buffer string) []string {
	if node == nil {
		return []string{}
	}

	var res []string

	tmpNode := node
	if tmpNode.pegRule == ruleID {
		for tmpNode != nil {
			res = append(res, buffer[tmpNode.begin:tmpNode.end])
			tmpNode = tmpNode.next
		}
	}

	return res
}

func (node *node32) validateBody(buffer string, scope *Scope, mergeScopes bool) error { // nolint // because
	if node == nil || node.up == nil || node.pegRule != ruleBODY {
		return nil
	}

	var tmpScope *Scope
	if mergeScopes {
		tmpScope = scope
	}
	bodyScope, statements, _, err := node.up.getParamsVars(buffer, tmpScope)
	if err != nil {
		return err
	}
	if bodyScope != scope {
		// new scope was created in getParamsVars
		bodyScope.up = scope
		bodyScope.currentFunction = scope.currentFunction
	}
	if statements.pegRule == ruleSTATEMENTS {
		statement := statements.up
		for statement != nil && statement.pegRule == ruleSTATEMENT {
			switch statement.up.pegRule {
			case ruleIF_STATEMENT:
				if err := statement.up.validateIfStatement(buffer, bodyScope); err != nil {
					return err
				}
			case ruleWHILE_STATEMENT:
				if err := statement.up.validateWhileStatement(buffer, scope); err != nil {
					return err
				}
			case ruleASSIGNMENT:
				if err := statement.up.validateAssignment(buffer, bodyScope); err != nil {
					return err
				}
			case ruleFUNC_CALL: // lebo viem modifikovat globalne premenne
				if err := statement.up.validateFuncCall(buffer, scope); err != nil {
					return err
				}
			}
			statement = statement.next
		}

		return nil
	}

	if statements.pegRule == ruleRETURN_CLAUSE {
		// validate return - check type of return value with function type
		retType, err := statements.up.up.getValueType(buffer, bodyScope)
		if err != nil {
			return err
		}
		funcType, isInScope := isVarInScope(bodyScope, bodyScope.currentFunction)
		if !isInScope {
			return fmt.Errorf("undefined function: %s", bodyScope.currentFunction)
		}
		if funcType != retType {
			return fmt.Errorf("cannot return %s value in function of type %s, in function %s", retType, funcType, bodyScope.currentFunction)
		}
	}

	return nil
}

func (node *node32) validateIfStatement(buffer string, scope *Scope) error {
	if node == nil || node.pegRule != ruleIF_STATEMENT {
		return nil
	}

	bodyNode, err := node.up.checkBoolExpression(buffer, scope)
	if err != nil {
		return fmt.Errorf("if statement: %s", err)
	}

	err = bodyNode.validateBody(buffer, scope, false)
	if err != nil {
		return fmt.Errorf("if statement: %s", err)
	}

	// checking else clause
	if bodyNode.next != nil {
		err = bodyNode.next.up.validateBody(buffer, scope, false)
		if err != nil {
			return fmt.Errorf("if statement: %s", err)
		}
	}

	return nil
}

func (node *node32) validateAssignment(buffer string, scope *Scope) error {
	if node == nil {
		return nil
	}

	if node.pegRule == ruleASSIGNMENT {
		varType, err := node.up.getValueType(buffer, scope)
		if err != nil {
			return fmt.Errorf("cannot determine variable type: %s", err)
		}

		value := node.up.next
		valueType, err := value.up.getValueType(buffer, scope)
		if err != nil {
			return fmt.Errorf("could not get type of value: %s", err)
		}
		if varType != valueType {
			return fmt.Errorf("cannot assign %s value to %s variable: %s", valueType, varType, buffer[node.begin:node.end])
		}

		if value.up.pegRule == ruleFUNC_CALL {
			return value.up.validateFuncCall(buffer, scope)
		}
	}

	return nil
}

func (node *node32) validateWhileStatement(buffer string, scope *Scope) error {
	if node == nil || node.pegRule != ruleWHILE_STATEMENT {
		return nil
	}

	bodyNode, err := node.up.checkBoolExpression(buffer, scope)
	if err != nil {
		return fmt.Errorf("while statement: %s", err)
	}

	err = bodyNode.validateBody(buffer, scope, false)
	if err != nil {
		return fmt.Errorf("while statement: %s", err)
	}

	return nil
}

func (node *node32) checkBoolExpression(buffer string, scope *Scope) (*node32, error) {
	if node == nil || node.pegRule != ruleBOOL_EXPR_VALUE {
		return node, nil
	}

	// check same types on every operation, whether it returns bool
	operandType, err := node.up.getValueType(buffer, scope)
	if err != nil {
		return node, err
	}
	tmp := node.next

	for tmp.pegRule != ruleBODY {
		if tmp.pegRule == ruleBOOL_EXPR_VALUE {
			valueType, _ := tmp.up.getValueType(buffer, scope)
			if operandType != valueType {
				return node, fmt.Errorf("cannot compare variables of different types: %s, %s", operandType, valueType)
			}
		}

		tmp = tmp.next
	}

	return tmp, nil
}

func (node *node32) getValueType(buffer string, scope *Scope) (VariableType, error) {
	switch node.pegRule {
	case ruleID:
		tmpScope := scope
		for tmpScope != nil {
			typ, isInScope := node.isVarInScope(tmpScope, buffer)
			if isInScope {
				return typ, nil
			}
			tmpScope = tmpScope.up
		}

		return "", fmt.Errorf("variable was not declared before used: %s", buffer[node.begin:node.end])
	case ruleTEXT:
		return String, nil
	case ruleINTEGER:
		return Integer, nil
	case ruleINT:
		return Integer, nil
	case ruleBOOLEAN:
		return Boolean, nil
	case ruleEXPRESSION:
		exprType, err := node.up.validateExpression(buffer, scope)
		if err != nil {
			return "", err
		}

		return exprType, nil
	case ruleFUNC_CALL:
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
		return id.variableType, true
	}

	return stringToVariableType("unknown"), false
}

func isVarInScope(scope *Scope, variable string) (VariableType, bool) {
	tmpScp := scope
	for tmpScp != nil {
		if id, ok := tmpScp.vars[variable]; ok {
			return id.variableType, true
		}
		tmpScp = tmpScp.up
	}

	return stringToVariableType("unknown"), false
}

func stringToID(scope *Scope, variable string) (*ID, bool) {
	tmpScp := scope
	for tmpScp != nil {
		if id, ok := tmpScp.vars[variable]; ok {
			return &id, true
		}
		tmpScp = tmpScp.up
	}

	return nil, false
}

func (node *node32) validateExpression(buffer string, scope *Scope) (VariableType, error) {
	if node == nil || node.pegRule != ruleEXPR_VALUE {
		return "", nil
	}

	typeOfExprValue, err := node.up.getValueType(buffer, scope)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	tmpNode := node.next
	for tmpNode != nil {
		if tmpNode.pegRule == ruleOP { // nolint
			op := buffer[tmpNode.begin:tmpNode.end]

			if typeOfExprValue == Boolean { // ID/func_call can be of type bool
				return "", fmt.Errorf("operation %s is not defined on boolean value", op)
			}
			if typeOfExprValue == String && (op == "-" || op == "*" || op == "/") {
				return "", fmt.Errorf("operation %s is not defined on string value", op)
			}
		} else if tmpNode.pegRule == ruleEXPR_VALUE {
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
	// FUNC_CALL / TEXT / INTEGER / ID -> (ID/ func_call moze mat typ bool)
	// +: b+b
	// *: s*s, b*b
	// -: s-s, b-b
	// /: s/s, b/b, (!) i/i- > check division by 0 when generating + vysledny typ
	return typeOfExprValue, nil
}

func (node *node32) validateFuncCall(buffer string, scope *Scope) error {
	if node == nil || node.pegRule != ruleFUNC_CALL {
		return nil
	}

	funcName := buffer[node.up.begin:node.up.end]
	id, found := stringToID(scope, funcName)
	if !found {
		return fmt.Errorf("function %s not defined before used", funcName)
	}

	fun, ok := id.additionalInfo.(Func)
	if !ok {
		return fmt.Errorf("calling variable in function call: %s", funcName)
	}

	ids := node.up.next.getVarIDs(buffer)
	if len(ids) != len(fun.paramsOrder) {
		return fmt.Errorf("inconsistent number of parameters in function call: %s", funcName)
	}

	for i, varType := range fun.paramsOrder {
		typ, exists := isVarInScope(scope, ids[i])
		if !exists {
			return fmt.Errorf("variable in function call %s does not exist: %s", funcName, ids[i])
		}
		if typ != varType {
			return fmt.Errorf("variable type does not match function parameter type: %s != %s", typ, varType)
		}
	}

	return nil
}