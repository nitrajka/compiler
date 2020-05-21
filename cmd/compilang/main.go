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
	string [raw, a, b, c]
	int [i, j, k, p]
	bool [s]
endglobals

func fibonacci(int [n]): int {
	int [res1, res2];
	if n <= 0 {; print("invalid input") }
    else {;
		if n == 1 {; return 0}
		if n == 2 {; return 1}
		k = n-1
		p = n-2
		res1 = call fibonacci(k)
		res2 = call fibonacci(p)
		return res1 + res2
    }

}

func fibonacci(int [n]): int {;
	print("zbytocna funkcia")
}

func emptyfunction(): void {;}

main
	{ string [a] ;
		if a==b {;}
		if 1 == -1 {;}
		if "ahoj" == "cau" {;} else {;}
		if a == b == c {;}
		if true == false {;}
		while a == b {;}
		newvar = -3
		k = var p
		z = "ahoj"
		msg = true
		print(a)
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