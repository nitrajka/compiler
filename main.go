package main

func main() {
	la := NewLexicalAnalyzer([]rune(";"))
	compiler := NewCompiler(la)
	compiler.Compile()
}