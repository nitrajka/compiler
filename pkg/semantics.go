package lexicalanalyser

import (
	"fmt"
	"strings"
)

type SemanticsError struct {
	buffer  string
	node    *node32
	context []string
}

func (e *SemanticsError) Error() string {
	line := strings.Count(e.buffer[:e.node.begin], "\n") + 1
	tmp := fmt.Sprintf("semantics error on line %d", line)
	msgs := append(e.context, tmp)

	for i := 0; i < len(msgs)/2; i++ {
		msgs[i], msgs[len(msgs)-i-1] = msgs[len(msgs)-1-i], msgs[i]
	}

	return strings.Join(msgs, "\n")
}

func AddErrorContext(err error, msg string) error {
	if se, ok := err.(*SemanticsError); ok {
		return &SemanticsError{buffer: se.buffer, node: se.node, context: append(se.context, msg)}
	}

	return fmt.Errorf("%s: %w", msg, err)
}

func NewSemanticsError(buffer string, node *node32) *SemanticsError {
	return &SemanticsError{buffer: buffer, node: node}
}

func NewSemanticsErrorf(buffer string, node *node32, format string, values ...interface{}) error {
	return AddErrorContext(&SemanticsError{buffer: buffer, node: node}, fmt.Sprintf(format, values...))
}

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

func (node *node32) CheckSemantics(buffer string) error {
	globalScope, tmpNode, _, err := node.up.getParamsVars(buffer, nil)
	if err != nil {
		return err
	}

	if tmpNode.pegRule == ruleFUNCTIONS {
		functionNode := tmpNode.up

		// declare all functions into global scope
		for functionNode != nil {
			id, err2 := functionNode.declareFunction(buffer, globalScope)
			if err2 != nil {
				return err2
			}

			globalScope.vars[buffer[functionNode.up.begin:functionNode.up.end]] = id
			functionNode = functionNode.next
		}

		// define all previously declared functions
		functionNode = tmpNode.up
		for functionNode != nil {
			id, err2 := functionNode.validateFunction(buffer, globalScope)
			if err2 != nil {
				return err2
			}

			globalScope.vars[buffer[functionNode.up.begin:functionNode.up.end]] = id // additionalInfo changes
			functionNode = functionNode.next
		}

		tmpNode = tmpNode.next
	}

	// validate main body
	globalScope.currentFunction = "main"
	globalScope.vars["main"] = ID{variableKind: Function, variableType: Void, name: "string"} // todo: need additionalInfo??
	err = tmpNode.validateBody(buffer, globalScope, false)

	if err != nil {
		return AddErrorContext(err, "main function is invalid")
	}

	return nil
}

func (node *node32) declareFunction(buffer string, scope *Scope) (ID, error) {
	funcName := buffer[node.up.begin:node.up.end]

	_, isInScope := node.up.isVarInScope(scope, buffer)
	if isInScope {
		return ID{}, NewSemanticsErrorf(buffer, node, "function name must be unique: %s", funcName)
	}

	functionScope, nodeType, paramsOrder, err := node.up.next.getParamsVars(buffer, nil)
	if err != nil {
		return ID{}, err // "function params vars: %s", err)
	}

	functionScope.up = scope
	functionScope.currentFunction = funcName

	funcType := buffer[nodeType.begin:nodeType.end]

	return ID{
		variableKind:   Function,
		variableType:   stringToVariableType(funcType),
		name:           funcName,
		additionalInfo: Func{paramsVars: functionScope, paramsOrder: paramsOrder, bodyValid: false},
	}, nil
}

func (node *node32) validateFunction(buffer string, scope *Scope) (ID, error) {
	functionNode := node

	functionScope, functionTypeNode, paramsOrder, err := functionNode.up.next.getParamsVars(buffer, nil)
	if err != nil {
		return ID{}, AddErrorContext(err, "function params vars")
	}

	functionScope.up = scope
	functionScope.currentFunction = buffer[functionNode.up.begin:functionNode.up.end]
	body := functionTypeNode.next

	if err := body.validateBody(buffer, functionScope, true); err != nil {
		return ID{}, err
	}
	//todo: validate print statement
	//todo: void global variables unnecessary (check and throw error?) + functions of void type cannot assign
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
				return nil, node, nil, NewSemanticsErrorf(buffer, tmpNode,
					"cannot define same variable ID more times in the same scope: %s", varID)
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

func (node *node32) validateBody(buffer string, scope *Scope, functionBody bool) error { // nolint // because
	if node == nil || node.pegRule != ruleBODY {
		return nil
	}

	id, exists := stringToID(scope, scope.currentFunction)
	if !exists || id.variableKind != Function {
		return NewSemanticsErrorf(buffer, node, "function %s does not exist", scope.currentFunction)
	}
	if node.up == nil {
		if functionBody && id.variableType != Void {
			return NewSemanticsErrorf(buffer, node, "function %s does not have return statement at the end of its body", scope.currentFunction)
		}
		return nil
	}

	var tmpScope *Scope
	if functionBody { // merge scopes of params and body declarations if in function body
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
				if err := statement.up.validateWhileStatement(buffer, bodyScope); err != nil {
					return err
				}
			case ruleASSIGNMENT:
				if err := statement.up.validateAssignment(buffer, bodyScope); err != nil {
					return err
				}
			case ruleFUNC_CALL: // lebo viem modifikovat globalne premenne
				if err := statement.up.validateFuncCall(buffer, bodyScope); err != nil {
					return err
				}
			}
			statement = statement.next
		}
		statements = statements.next
	}

	id, exists = stringToID(scope, scope.currentFunction)
	if !exists || id.variableKind != Function {
		return NewSemanticsErrorf(buffer, statements, "function %s does not exist", scope.currentFunction)
	}
	if functionBody && id.variableType != Void && statements == nil {
		return NewSemanticsErrorf(buffer, node, "function %s does not have return at the end of body", scope.currentFunction)
	}
	if statements != nil && statements.pegRule == ruleRETURN_CLAUSE {
		// validate return - check type of return value with function type
		var retType VariableType
		if statements.up.up == nil { // func_call's value does not have child if Void
			retType = Void
		} else {
			retType, err = statements.up.up.getValueType(buffer, bodyScope)
			if err != nil {
				return err
			}
		}

		funcType, isInScope := isVarInScope(bodyScope, bodyScope.currentFunction)
		if !isInScope {
			return NewSemanticsErrorf(buffer, node, "undefined function: %s", bodyScope.currentFunction)
		}

		if funcType != retType {
			return NewSemanticsErrorf(buffer, node,
				"cannot return %s value in function of type %s, in function %s", retType, funcType, bodyScope.currentFunction)
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
		return AddErrorContext(err, "if statement")
	}

	err = bodyNode.validateBody(buffer, scope, false)
	if err != nil {
		return AddErrorContext(err, "if statement")
	}

	// checking else clause
	if bodyNode.next != nil {
		err = bodyNode.next.up.validateBody(buffer, scope, false)
		if err != nil {
			return AddErrorContext(err, "if statement")
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
			return AddErrorContext(err, "cannot determine variable type")
		}

		value := node.up.next

		valueType, err := value.up.getValueType(buffer, scope)
		if err != nil {
			return AddErrorContext(err, "could not get type of value")
		}

		if varType != valueType {
			return NewSemanticsErrorf(buffer, node, "cannot assign %s value to %s variable: %s", valueType, varType, buffer[node.begin:node.end])
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
		return AddErrorContext(err, "while statement")
	}

	err = bodyNode.validateBody(buffer, scope, false)
	if err != nil {
		return AddErrorContext(err, "while statement")
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
				return node, NewSemanticsErrorf(buffer, node, "cannot compare variables of different types: %s, %s", operandType, valueType)
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

		return "", NewSemanticsErrorf(buffer, node, "variable was not declared before used: %s", buffer[node.begin:node.end])
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
		return "", NewSemanticsErrorf(buffer, node, "variable was not declared: %s", buffer[node.begin:node.end])
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
		return "", err
	}

	tmpNode := node.next
	for tmpNode != nil {
		if tmpNode.pegRule == ruleOP { // nolint
			op := buffer[tmpNode.begin:tmpNode.end]

			if typeOfExprValue == Boolean { // ID/func_call can be of type bool
				return "", NewSemanticsErrorf(buffer, node, "operation %s is not defined on boolean value", op)
			}

			if typeOfExprValue == String && (op == "-" || op == "*" || op == "/" || op == "%") {
				return "", NewSemanticsErrorf(buffer, node, "operation %s is not defined on string value", op)
			}
		} else if tmpNode.pegRule == ruleEXPR_VALUE {
			nextType, err := tmpNode.up.getValueType(buffer, scope)

			if err != nil {
				return "", err
			}

			if nextType != typeOfExprValue {
				return "", NewSemanticsErrorf(buffer, node, "cannot operate on values of different types: %s, %s", typeOfExprValue, nextType)
			}
		}

		tmpNode = tmpNode.next
	}
	// FUNC_CALL / TEXT / INTEGER / ID -> (ID/ func_call moze mat typ bool)
	// %: s%s, b%b
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
		return NewSemanticsErrorf(buffer, node, "function %s not defined before used", funcName)
	}

	fun, ok := id.additionalInfo.(Func)
	if !ok {
		return NewSemanticsErrorf(buffer, node, "calling variable in function call: %s", funcName)
	}

	ids := node.up.next.getVarIDs(buffer)
	if len(ids) != len(fun.paramsOrder) {
		return NewSemanticsErrorf(buffer, node, "inconsistent number of parameters in function call: %s", funcName)
	}

	for i, varType := range fun.paramsOrder {
		typ, exists := isVarInScope(scope, ids[i])
		if !exists {
			return NewSemanticsErrorf(buffer, node, "variable in function call %s does not exist: %s", funcName, ids[i])
		}

		if typ != varType {
			return NewSemanticsErrorf(buffer, node, "variable type does not match function parameter type: %s != %s", typ, varType)
		}
	}

	return nil
}