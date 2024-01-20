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
	readFile, err := os.ReadFile(fileName + ".py")
	if err != nil {
		log.Fatal("[main (compile_file)] File error")
		return
	}

	output := compile(readFile)

	// Write to the file and close it
	f, err := os.Create("../Output/" + fileName + ".go")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(output)
	err = f.Sync()
	if err != nil {
		fmt.Println(err)
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func compile(input []byte) string {
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
	//fmt.Println(ast)
	//fmt.Println(ast.stringify())

	// Optimize

	// Emit
	emitter := Emitter{}
	emitSource, err := emitter.emit(ast)
	if err != nil {
		log.Fatal(err)
	}
	// Final code
	//fmt.Println(emitSource)
	return "package main\nfunc main() {\n" + emitSource + "}"
}
