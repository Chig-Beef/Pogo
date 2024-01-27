package main

import (
	"errors"
	"strconv"
	"strings"
)

type Emitter struct {
}

func (e *Emitter) emit(ast Structure) (string, error) {
	output := ""
	if ast.code == structureCode["ILLEGAL"] {
		return output, errors.New("[Emit (emit)] ILLEGAL structure found in final code" + " on line " + strconv.Itoa(ast.line))
	}

	// Override for loops
	if ast.code == structureCode["ST_FOR"] {
		identifier := ast.children[1].text
		output += "\n" + ast.children[0].text + " " + identifier + " := " + ast.children[5].text + ";"
		output += identifier + " < " + ast.children[7].text + ";"
		output += identifier + "++"
		temp, err := e.emit(ast.children[10])
		if err != nil {
			return output, err
		}
		output += " " + temp
		return output, nil
	}

	if ast.code == structureCode["ST_WHILE"] {
		output += "for "
		temp, err := e.emit(ast.children[1])
		if err != nil {
			return output, err
		}
		output += temp
		temp, err = e.emit(ast.children[4])
		if err != nil {
			return output, err
		}
		output += temp
		return output, nil
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

	// Keywords
	structureCode["K_IF"]:    "\nif",
	structureCode["K_ELIF"]:  "else if",
	structureCode["K_WHILE"]: "\nfor",
	structureCode["K_DEF"]:   "\nfunc",

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
	structureCode["FUNC_NAME"],

	// Keywords
	structureCode["K_ELSE"],
	structureCode["K_RETURN"],

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
