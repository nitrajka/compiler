package lexicalanalyser

import (
	"fmt"
	"os"
)

func (t *tokens32) WalkAndDeleteUnwanted(buffer string) {
	node := t.AST()
	node.deleteWhitespace()
	node.parsePARAMS_VARS()
	//node.PrettyPrint(os.Stdout, buffer)
	//node.parseFunctions()
	node.parseStatements()
	node.parseIfWhileStatement(buffer)
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

func (node *node32) parsePARAMS_VARS() {
	if node == nil {
		return
	}

	if rul3s[node.pegRule] == rul3s[rulePARAMS_VARS] {
		//typ := strconv.Quote(string(([]rune(buffer)[node.up.begin:node.up.end])))
		ids := node.up.next.parseVAR_LIST()

		if node.up.next.next != nil {
			tmp := node.next
			node.next = node.up.next.next
			node.next.next = tmp
		}

		//fmt.Println(ids)
		if len(ids) > 0 {
			node.up = &node32{token32: token32{pegRule: ruleTYPE, begin: node.up.begin, end: node.up.end}}
			node.up.up = &node32{token32: token32{pegRule: ruleID, begin: ids[0].begin, end: ids[0].end}}
			node1 := &node.up.up
			for i := 1; i < len(ids); i++ {
				(*node1).next = &node32{token32: token32{pegRule: ruleID, begin: ids[i].begin, end: ids[i].end}}
				node1 = &(*node1).next
			}
		}

	} else {
		if node.up != nil {
			node.up.parsePARAMS_VARS()
		}
	}
	if node.next != nil {
		node.next.parsePARAMS_VARS()
	}
}

type ID struct {
	begin, end uint32
}

func (node *node32) parseVAR_LIST() []ID {
	if node == nil {
		return nil
	}
	rule := rul3s[node.pegRule]
	var ids []ID

	if rule == rul3s[ruleVAR_LIST] {
		node1 := node.up
		ruleNode := rul3s[node1.pegRule]
		if ruleNode == rul3s[ruleID] {
			ids = append(ids, ID{begin: node1.token32.begin, end: node1.token32.end})
		}
		node2 := node1.next
		for node2 != nil {
			ids = append(ids, node2.parseVAR_LIST()...)
			node2 = node2.next
		}
		if node1.up != nil {
			node1.up.parseVAR_LIST()
		}
	} else if rule == rul3s[ruleVAR_LIST1] {
		return  node.up.parseVAR_LIST()
	}

	return ids
}

//func (node *node32) parseFunctions() {
//	if node == nil {
//		return
//	}
//
//}

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

func (node *node32) parseIfWhileStatement(buffer string) {
	if node == nil {
		return
	}
	rule := rul3s[node.pegRule]
	if rule == rul3s[ruleIF_STATEMENT] || rule == rul3s[ruleWHILE_STATEMENT] {
		tmp := node.up.next
		node.up = node.up.up
		tmpNode := node.up
		for tmpNode.next != nil {
			tmpNode = tmpNode.next
		}
		tmpNode.next = tmp
	} else if rule == rul3s[ruleBOOL_OP] {
		tmp := node.next.next
		node.next = node.next.up
		tmpNode := node.next
		for tmpNode.next != nil {
			tmpNode = tmpNode.next
		}
		tmpNode.next = tmp
	} else if rule == rul3s[ruleBOOL_EXPR_VALUE] {
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
		node.up.parseIfWhileStatement(buffer)
	}
	if node.next != nil {
		node.next.parseIfWhileStatement(buffer)
	}
}

func (node *node32) parseStrings() {
	if node == nil {
		return
	}
	rule := rul3s[node.pegRule]
	if rule == rul3s[ruleTEXT] {
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
	rule := rul3s[node.pegRule]
	if rule == rul3s[ruleINTEGER] {
		node.up = nil
	}
	if node.up != nil {
		node.up.parseIntegers()
	}
	if node.next != nil {
		node.next.parseIntegers()
	}
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

	if node.up != nil {
		node.up.parseAssignment()
	}

	if node.next != nil {
		node.next.parseAssignment()
	}
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

	if node.up != nil {
		node.up.parseExpression()
	}

	if node.next != nil {
		node.next.parseExpression()
	}
}