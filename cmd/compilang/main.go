package main

import (
	"fmt"
	"os"

	"github.com/nitrajka/compiler/pkg"

	"github.com/nitrajka/compiler/pkg/lexicalanalyser"
)

func main() {
	//todo: nacitat subor, obsah skompilovat,
	//todo: meno suboru cez argument
	//todo: readme - ako to spustit, skompilovat
	la := lexicalanalyser.NewLexicalAnalyzer("globals endglobals main { ; return void } endmain")
	compiler := pkg.NewCompiler(la)
	program, err := compiler.Compile()
	if err != nil {
		fmt.Printf("compilation failed: %v", err)
		os.Exit(1)
	}

	fmt.Println(program)
}