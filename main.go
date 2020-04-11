package main

import (
	"fmt"
	"github.com/nitrajka/compiler/lexicalAnalyzator"
	"os"
)

func main() {
	la := lexicalAnalyzator.NewLexicalAnalyzer("globals endglobals main { ; return void } endmain")
	compiler := NewCompiler(la)
	program, err := compiler.Compile()
	if err != nil {
		fmt.Printf("compilation failed: %v", err)
		os.Exit(1)
	}

	fmt.Println(program)
}