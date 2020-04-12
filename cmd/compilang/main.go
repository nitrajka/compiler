package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nitrajka/compiler/pkg"

	"github.com/nitrajka/compiler/pkg/lexicalanalyser"
)

func main() {
	flag.Parse()
	programFileName := flag.Arg(0)
	if programFileName == "" {
		exit(fmt.Sprintf("provide program file to compile"))
	}

	content, err := ioutil.ReadFile(programFileName)
	if err != nil {
		exit(fmt.Sprintf("reading of file failed %v", err))
	}

	//todo: readme - ako to spustit, skompilovat
	la := lexicalanalyser.NewLexicalAnalyzer(string(content))
	compiler := pkg.NewCompiler(la)
	program, err := compiler.Compile()
	if err != nil {
		exit(fmt.Sprintf("compilation failed: %v", err))
	}

	fmt.Println(program)
}

func exit(err string) {
	fmt.Println(err)
	os.Exit(1)
}