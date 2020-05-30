package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	lexicalanalyser "github.com/nitrajka/compiler/pkg"
)

func main() {
	output := flag.String("o", "a.go", "Output for compiler")
	flag.Parse()
	programFileName := flag.Arg(0)

	if programFileName == "" {
		exit("provide program file to compile")
	}

	content, err := ioutil.ReadFile(programFileName)
	if err != nil {
		exit(fmt.Sprintf("reading of file failed %v", err))
	}

	parser := &lexicalanalyser.MyParser{Buffer: string(content), Pretty: true}
	err = parser.Init()

	if err != nil {
		exit(err.Error())
	}

	if err := parser.Parse(1); err != nil {
		exit(err.Error())
	}

	root := parser.ParseAST(string(content))
	//fmt.Println("OK: ast preparsed")
	err2 := root.CheckSemantics(string(content))
	if err2 != nil {
		fmt.Println("FAIL: code semantics")
		exit(err2.Error())
	}
	//fmt.Println("OK: code semantics")

	err3 := root.Generate(string(content), *output)
	if err3 != nil {
		fmt.Println("FAIL: code generating")
		exit(err3.Error())
	}
	//fmt.Println("OK: code generated")
}

func exit(err string) {
	fmt.Println(err)
	os.Exit(1)
}