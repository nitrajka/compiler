package main

import (
	"fmt"
	"os"

	"github.com/nitrajka/compiler/pkg"
)

func main() {
	//flag.Parse()
	//programFileName := flag.Arg(0)
	//if programFileName == "" {
	//	exit(fmt.Sprintf("provide program file to compile"))
	//}
	//
	//content, err := ioutil.ReadFile(programFileName)
	//if err != nil {
	//	exit(fmt.Sprintf("reading of file failed %v", err))
	//}

	//	if n < 0 {;
	//	} else {
	//		if n==1 {;} else if n==2 {; return 1} else {; return fibonacci(n-1)+fibonacci(n-2) }
	//
	//	}
	content := `globals
	string [raw]
	int [i, j, k]
endglobals
main
	{;
		if a==b {;}
		if arr a[1] == arr b[1] {;}
		return void
	}
endmain
`
	parser := &lexicalanalyser.MyParser{Buffer: string(content), Pretty: true}
	err := parser.Init()
	if err != nil {
		exit(err.Error())
	}
	if err := parser.Parse(1); err != nil {
		exit(err.Error())
	}

	//ast := parser.AST()
	//parser.PrettyPrintSyntaxTree(content)
	parser.WalkAndDeleteUnwanted(content)
	//parser.Walk()
}

func exit(err string) {
	fmt.Println(err)
	os.Exit(1)
}