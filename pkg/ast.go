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

func isWhitespace(rule pegRule) bool {
	return rule == ruleWHITESPACE_AT_LEAST_ONE || rule == ruleWHITESPACE_ANY ||
		rule == ruleWHITESPACEC || rule == ruleFIRST_CHAR || rule == ruleID1 ||
		rule == ruleJUST_SPACES || rule == ruleAT_LEAST_ONE_SPACE
}

func (node *node32) deleteWhitespace() {
	if node == nil {
		return
	}

	for node.up != nil && isWhitespace(node.up.pegRule) {
		node.up = node.up.next
	}

	node.up.deleteWhitespace()

	node1 := node.up
	for node1 != nil && node1.next != nil {
		if isWhitespace(node1.next.pegRule) {
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

	if node.pegRule == rulePARAMS_VARS {
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
	} else { // must have else if
		node.up.parseParamsVars()
	}

	node.next.parseParamsVars()
}

func (node *node32) parseFuncCall() {
	if node == nil {
		return
	}

	if node.pegRule == ruleFUNC_CALL {
		if node.up.next != nil {
			node.up.next.parseVarList()
			node.up.next = node.up.next.up
		}
	}

	node.up.parseFuncCall()
	node.next.parseFuncCall()
}

func (node *node32) parseVarList() {
	if node == nil {
		return
	}

	if node.pegRule == ruleVAR_LIST {
		if node.up.next != nil {
			node.up.next = node.up.next.up.up
		}

		node.up.next.parseVarList()
	} else if node.pegRule == ruleID {
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

	if node.pegRule == ruleSTATEMENT && node.next != nil {
		node.next = node.next.up
	}

	node.next.parseStatements()
	node.up.parseStatements()
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
		if node.next.next != nil && node.next.next.pegRule != ruleELSECLAUSE {
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

	node.up.parseIfWhileStatement()
	node.next.parseIfWhileStatement()
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

	node.up.parseStrings()
	node.next.parseStrings()
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

	if node.pegRule == ruleASSIGNMENT {
		value := node.up.next
		node.up = node.up.up

		if node.up.next != nil { // handle indexing of an array
			node.up.up = node.up.next.up
			node.up.next = nil
		}

		node.up.next = value

		if value.up.next != nil && value.up.next.pegRule == ruleINDEXED {
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

	if node.pegRule == ruleEXPR_VALUE || node.pegRule == ruleOP {
		if node.next != nil {
			node.next = node.next.up
		}
	}

	node.up.parseExpression()
	node.next.parseExpression()
}
