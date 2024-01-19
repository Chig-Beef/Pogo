package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 2 {
		compile_file(args[1])
	}
}

func compile_file(fileName string) {
	readFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("[main (compile_file)] File error")
		return
	}

	compile(readFile)
}

func compile(input []byte) {
	//fmt.Println(input)

	// Lex
	lexer := Lexer{}
	lexSource := lexer.lex(input)
	//fmt.Println(lexSource)

	// Parse
	parser := Parser{}
	pS := parser.replaceIndents(lexSource)
	//fmt.Println(pS)
	ast := parser.parse(pS)

	// Optimize

	// Emit

	fmt.Println(ast.stringify())
}
