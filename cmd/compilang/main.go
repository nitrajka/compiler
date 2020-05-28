package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	lexicalanalyser "github.com/nitrajka/compiler/pkg"
)

func main() {
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

	parser.WalkAndDeleteUnwanted(string(content))
}

func exit(err string) {
	fmt.Println(err)
	os.Exit(1)
}