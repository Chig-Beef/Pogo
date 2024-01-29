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
	fmt.Println(ast.stringify())

	main_func := Structure{
		structureCode["ST_FUNCTION"],
		"ST_FUNCTION",
		-1,
		[]Structure{
			createStructure("K_DEF", "def", -1),
			createStructure("FUNC_NAME", "main", -1),
			createStructure("L_PAREN", "(", -1),
			createStructure("R_PAREN", ")", -1),
			createStructure("COLON", ":", -1),
			{structureCode["BLOCK"], "", -1, append(ast.children, createStructure("ANTI_COLON", ":", -1))},
		},
	}

	ast.children = []Structure{main_func}
	ast.children = append(parser.functions, ast.children...)

	//fmt.Println(ast.stringify())

	// Analyze
	analyzer := Analyzer{}
	err := analyzer.analyze(ast, []Variable{}, []Function{{"print", []string{"any"}, "None"}})
	if err != nil {
		log.Fatal(err)
	}

	// Optimize

	// Emit
	emitter := Emitter{}
	emitSource, err := emitter.emit(ast)
	if err != nil {
		log.Fatal(err)
	}
	// Final code
	//fmt.Println(emitSource)
	return "package main\n" + emitSource
}
