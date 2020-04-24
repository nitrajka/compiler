package lexicalanalyser

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func (t *tokens32) WalkAndDeleteUnwanted(buffer string) {
	t.AST().WalkAndDeleteWhitespace(os.Stdout, true, buffer)
}

func (node *node32) WalkAndDeleteWhitespace(w io.Writer, pretty bool, buffer string) { //func(n *node32) *node32
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))

			if node.up != nil { // has child
				childRule := rul3s[node.up.pegRule]
				for node.isWhitespace(childRule) && node.up != nil {
					node.up = node.up.next
					if node.up != nil {
						childRule = rul3s[node.up.pegRule]
					}
				}

				if node.up != nil {
					node1 := node.up
					for node1.next != nil {
						ruleNext := rul3s[node1.next.pegRule]
						if node.isWhitespace(ruleNext) {
							node1.next = node1.next.next
						} else {
							node1 = node1.next
						}
					}
				}
			}

			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) isWhitespace(rule string) bool {
	return rule == rul3s[ruleWHITESPACE_AT_LEAST_ONE] || rule == rul3s[ruleWHITESPACE_ANY] ||
		rule == rul3s[ruleWHITESPACEC] || rule == rul3s[ruleFIRST_CHAR] ||
		rule == rul3s[ruleJUST_SPACES] || rule == rul3s[ruleAT_LEAST_ONE_SPACE]
}

func (node *node32) parsePARAMS_VARS(buffer string) {
	rule := rul3s[node.up.pegRule]
	var types []string
	if rule == rul3s[rulePARAMS_VARS] {
		quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
		if node.up != nil { // at least 1 variable
			types = append(types, quote)
			ids := node.up.parseVAR_LIST(buffer)
			if len(ids) > 0 {
				node.up = &node32{token32: token32{pegRule: ruleID, begin: ids[0].begin, end: ids[0].end}}
				node1 := node.up
				for i := 1; i < len(ids); i++ {
					node1.next = &node32{token32:token32{pegRule: ruleID, begin: ids[i].begin, end: ids[i].end}}
					node1 = node1.next
				}
			}

		}
	}
}

type ID struct {
	begin, end uint32
}

func (node *node32) parseVAR_LIST(buffer string) []ID {
	rule := rul3s[node.up.pegRule]
	var ids []ID

	if rule == rul3s[ruleVAR_LIST] {
		queue := []*node32{node}
		for len(queue) > 0 {
			node1 := queue[0]
			queue = queue[1:]
			ruleNode := rul3s[node1.pegRule]
			if ruleNode == rul3s[ruleID] {
				ids = append(ids, ID{begin:node1.begin, end: node1.end})
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