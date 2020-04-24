package lexicalanalyser

import (
	"io/ioutil"
	"testing"
)

func TestGrammarValid(t *testing.T) {
	testcases := map[string]bool {
		"cases/case1.txt": true,
		"cases/case2.txt": true,
		"cases/case3.txt": false,
		"cases/case4.txt": false,
		"cases/fibonacci.txt": true,
	}

	parser := &MyParser{Pretty: true}
	err := parser.Init()
	if err != nil {
		t.Errorf(err.Error())
	}

	for fileName, isValid := range testcases {

		t.Run("test program valid", func(t *testing.T) {
			content, err := ioutil.ReadFile(fileName)
			if err != nil {
				t.Errorf(err.Error())
			}

			parser.Buffer = string(content)
			parser.Reset()
			if err := parser.Parse(1); err != nil {
				if isValid {
					t.Errorf("valid program parsed as invalid: %s", err.Error())
				}

			} else if !isValid {
				t.Errorf("invalid program parsed as valid: %s", fileName)
			}
		})
	}
}

func TestSemanticsValid(t *testing.T) {

}