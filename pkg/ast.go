package lexicalanalyser

import (
	"fmt"
	"os"
	"strconv"
)

func (t *tokens32) WalkAndDeleteUnwanted(buffer string) {
	node := t.AST()
	node.deleteWhitespace()
	node.PrettyPrint(os.Stdout, buffer)
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

func (node *node32) parsePARAMS_VARS(buffer string) {
	if node.up == nil {
		return
	}
	rule := rul3s[node.up.pegRule]
	var types []string
	if rule == rul3s[rulePARAMS_VARS] {
		quote := strconv.Quote(string(([]rune(buffer)[node.up.up.begin:node.up.up.end])))
		types = append(types, quote)
		ids := node.up.up.next.next.parseVAR_LIST(buffer)

		fmt.Println(ids)
		if len(ids) > 0 {
			node.up = &node32{token32: token32{pegRule: ruleID, begin: ids[0].begin, end: ids[0].end}}
			node1 := &node.up
			for i := 1; i < len(ids); i++ {
				(*node1).next = &node32{token32: token32{pegRule: ruleID, begin: ids[i].begin, end: ids[i].end}}
				node1 = &(*node1).next
			}
		}
	}
}

type ID struct {
	begin, end uint32
}

func (node *node32) parseVAR_LIST(buffer string) []ID {
	rule := rul3s[node.pegRule]
	var ids []ID

	if rule == rul3s[ruleVAR_LIST] {
		queue := []*node32{node.up}
		for len(queue) > 0 {
			node1 := queue[0]
			queue = queue[1:]
			ruleNode := rul3s[node1.pegRule]
			if ruleNode == rul3s[ruleID] {
				ids = append(ids, ID{begin: node1.begin, end: node1.end})
			}
			if node1.next != nil {
				node2 := node1
				for node2 != nil {
					queue = append(queue, node2)
					node2 = node2.next
				}
			}
			if node1.up != nil {
				queue = append(queue, node1.up)
			}
		}
	}

	return ids
}
