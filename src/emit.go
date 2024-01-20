package main

import (
	"errors"
	"strings"
)

type Emitter struct {
}

func (e *Emitter) emit(ast Structure) (string, error) {
	output := ""
	if ast.code == structureCode["ILLEGAL"] {
		return output, errors.New("[Emit (emit)] ILLEGAL structure found in final code")
	}

	// Own text
	val, exists := translation[ast.code]
	if !exists {
		found := false
		for i := 0; i < len(directs); i++ {
			if directs[i] == ast.code {
				found = true
				output += ast.text
				break
			}
		}

		if !found {
			// Bools
			if ast.code == structureCode["L_BOOL"] {
				output += strings.ToLower(ast.text)
			}
		}
	}
	output += val

	// Children's text
	for i := 0; i < len(ast.children); i++ {
		temp, err := e.emit(ast.children[i])
		if err != nil {
			return output, err
		}
		output += " " + temp
	}
	return output, nil
}

var translation map[int]string = map[int]string{
	// Statements
	structureCode["ST_DECLARATION"]: "var",

	// Other
	structureCode["BLOCK"]:      "{",
	structureCode["NEWLINE"]:    "\n",
	structureCode["ANTI_COLON"]: "}",

	// In-built functions
	structureCode["IB_PRINT"]: "println",

	// Bool operands
	structureCode["BO_NOT"]: "!",
	structureCode["BO_AND"]: "&&",
	structureCode["BO_OR"]:  "||",
}

var directs []int = []int{
	// Other
	structureCode["IDENTIFIER"],
	structureCode["L_PAREN"],
	structureCode["R_PAREN"],
	structureCode["L_BLOCK"],
	structureCode["R_BLOCK"],
	structureCode["L_SQUIRLY"],
	structureCode["R_SQUIRLY"],
	structureCode["SEP"],
	structureCode["ASSIGN"],
	structureCode["COMMENT_ONE"],

	// Keywords
	structureCode["K_FOR"],
	structureCode["K_IF"],
	structureCode["K_ELIF"],
	structureCode["K_ELSE"],

	// Math operands
	structureCode["MO_PLUS"],
	structureCode["MO_SUB"],
	structureCode["MO_MUL"],
	structureCode["MO_DIV"],
	structureCode["MO_MODULO"],

	// Literal
	structureCode["L_INT"],
	structureCode["L_STRING"],

	// Comparison operands
	structureCode["CO_EQUALS"],
	structureCode["CO_NOT_EQUALS"],
	structureCode["CO_GT"],
	structureCode["CO_GT_EQUALS"],
	structureCode["CO_LT"],
	structureCode["CO_LT_EQUALS"],
}
