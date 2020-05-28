package lexicalanalyser

import (
	"fmt"
	"os"
)

func (t *tokens32) WalkAndDeleteUnwanted(buffer string) {
	node := t.AST()
	node.deleteWhitespace()
	node.parseParamsVars()
	node.parseFuncCall()
	node.parseStatements()
	node.parseIfWhileStatement()
	node.parseStrings()
	node.parseIntegers()
	node.parseAssignment()
	node.parseExpression()
	node.PrettyPrint(os.Stdout, buffer)

	err := node.checkSemantics(buffer)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("program ok")
	}
}

func isWhitespace(rule string) bool {
	return rule == rul3s[ruleWHITESPACE_AT_LEAST_ONE] || rule == rul3s[ruleWHITESPACE_ANY] ||
		rule == rul3s[ruleWHITESPACEC] || rule == rul3s[ruleFIRST_CHAR] || rule == rul3s[ruleID1] ||
		rule == rul3s[ruleJUST_SPACES] || rule == rul3s[ruleAT_LEAST_ONE_SPACE]
}

func (node *node32) deleteWhitespace() {
	for node.up != nil && isWhitespace(rul3s[node.up.pegRule]) {
		node.up = node.up.next
	}

	if node.up != nil {
		node.up.deleteWhitespace()
	}

	node1 := node.up
	for node1 != nil && node1.next != nil {
		if isWhitespace(rul3s[node1.next.pegRule]) {
			node1.next = node1.next.next
		} else {
			node1 = node1.next
			node1.deleteWhitespace()
		}
	}
}

func (node *node32) parseParamsVars() {
	if node == nil {
		return
	}

	if rul3s[node.pegRule] == rul3s[rulePARAMS_VARS] {
		node.up.next.parseVarList()

		if node.up.next.next != nil { // next params vars
			tmp := node.next
			node.next = node.up.next.next
			node.up.next.next = nil
			node.next.next = tmp
		}

		if node.up.next != nil {
			node.up.up = node.up.next.up
			node.up.next = nil
		}
	} else if node.up != nil { // must have else if
		node.up.parseParamsVars()
	}

	if node.next != nil {
		node.next.parseParamsVars()
	}
}

func (node *node32) parseFuncCall() {
	if node == nil {
		return
	}

	if rul3s[node.pegRule] == rul3s[ruleFUNC_CALL] {
		if node.up.next != nil {
			node.up.next.parseVarList()
			node.up.next = node.up.next.up
		}
	}

	if node.up != nil {
		node.up.parseFuncCall()
	}

	if node.next != nil {
		node.next.parseFuncCall()
	}
}

func (node *node32) parseVarList() {
	if node == nil {
		return
	}

	if rul3s[node.pegRule] == rul3s[ruleVAR_LIST] {
		if node.up.next != nil {
			node.up.next = node.up.next.up.up
		}
		node.up.next.parseVarList()
	} else if rul3s[node.pegRule] == rul3s[ruleID] {
		if node.next != nil {
			node.next = node.next.up.up
		}
		node.next.parseVarList()
	}
}

func (node *node32) parseStatements() {
	if node == nil {
		return
	}
	rule := rul3s[node.pegRule]
	if rule == rul3s[ruleSTATEMENT] {
		if node.next != nil {
			node.next = node.next.up
			node.next.parseStatements()
		}
	} else if node.next != nil {
		node.next.parseStatements()
	}
	if node.up != nil {
		node.up.parseStatements()
	}
}

func (node *node32) parseIfWhileStatement() {
	if node == nil {
		return
	}

	switch node.pegRule {
	case ruleIF_STATEMENT, ruleWHILE_STATEMENT:
		tmp := node.up.next
		node.up = node.up.up
		tmpNode := node.up

		for tmpNode.next != nil {
			tmpNode = tmpNode.next
		}

		tmpNode.next = tmp

	case ruleBOOL_OP:
		tmp := node.next.next
		node.next = node.next.up
		tmpNode := node.next

		for tmpNode.next != nil {
			tmpNode = tmpNode.next
		}

		tmpNode.next = tmp

	case ruleBOOL_EXPR_VALUE:
		if node.next.next != nil && rul3s[node.next.next.pegRule] != rul3s[ruleELSECLAUSE] {
			tmp := node.next.next
			node.next = node.next.up
			tmpNode := node.next

			if tmpNode != nil {
				for tmpNode.next != nil {
					tmpNode = tmpNode.next
				}
				tmpNode.next = tmp
			}
		}
	}

	if node.up != nil {
		node.up.parseIfWhileStatement()
	}

	if node.next != nil {
		node.next.parseIfWhileStatement()
	}
}

func (node *node32) parseStrings() {
	if node == nil {
		return
	}

	if node.pegRule == ruleTEXT {
		if node.up != nil { // if is empty string, has no STRING child
			node.begin = node.up.begin
			node.end = node.up.end
			node.token32 = token32{pegRule: ruleTEXT, begin: node.begin, end: node.end}
			node.up = nil
		}
	}

	if node.up != nil {
		node.up.parseStrings()
	}

	if node.next != nil {
		node.next.parseStrings()
	}
}

func (node *node32) parseIntegers() {
	if node == nil {
		return
	}

	if node.pegRule == ruleINTEGER {
		node.up = nil
	}

	node.up.parseIntegers()
	node.next.parseIntegers()
}

func (node *node32) parseAssignment() {
	if node == nil {
		return
	}

	if rul3s[node.pegRule] == rul3s[ruleASSIGNMENT] {
		value := node.up.next
		node.up = node.up.up

		if node.up.next != nil { // handle indexing of an array
			node.up.up = node.up.next.up
			node.up.next = nil
		}
		node.up.next = value

		if value.up.next != nil && rul3s[value.up.next.pegRule] == rul3s[ruleINDEXED] {
			value.up.up = value.up.next.up.up
			value.up.next = nil
		}
	}

	node.up.parseAssignment()
	node.next.parseAssignment()
}

func (node *node32) parseExpression() {
	if node == nil {
		return
	}

	if rul3s[node.pegRule] == rul3s[ruleEXPR_VALUE] || rul3s[node.pegRule] == rul3s[ruleOP] {
		if node.next != nil {
			node.next = node.next.up
		}
	}

	node.up.parseExpression()
	node.next.parseExpression()
}
