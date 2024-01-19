package main

import "strings"

type Structure struct {
	children []Structure
	code     int
	text     string
}

func (st Structure) stringify() string {
	text := ""
	for i := 0; i < len(st.children); i++ {
		text += "\n" + st.children[i].stringify()
	}
	var t string
	if st.text == "\n" {
		t = "NEWLINE"
	} else {
		t = st.text
	}
	return t + strings.ReplaceAll(text, "\n", "\n\t")
}

var structureCode map[string]int = map[string]int{
	// Not implemented
	"ILLEGAL": -1,

	"ST_IMPORT":     0,
	"IF_ELSE_BLOCK": 1,
	"BLOCK":         2,
	"ST_IF":         3,
	"ST_ELIF":       4,
	"ST_ELSE":       5,
	"CALL":          6,
	"ST_FOR":        7,
	"EXPRESSION":    8,
	"PROGRAM":       9,
	"STATEMENT":     10,
	"ASTERISK":      11,
	"K_IMPORT":      12,
	"K_FROM":        13,
	"K_FOR":         14,
	"K_IN":          15,
	"K_IF":          16,
	"K_ELIF":        17,
	"K_ELSE":        18,
	"IB_PRINT":      19,
	"IB_RANGE":      20,
	"BO_NOT":        21,
	"BO_AND":        22,
	"BO_OR":         23,
	"MO_PLUS":       24,
	"MO_SUB":        25,
	"MO_MUL":        26,
	"MO_DIV":        27,
	"MO_MODULO":     28,
	"IDENTIFIER":    29,
	"NEWLINE":       30,
	"INDENT":        31,
	"L_PAREN":       32,
	"R_PAREN":       33,
	"L_BLOCK":       34,
	"R_BLOCK":       35,
	"L_SQUIRLY":     36,
	"R_SQUIRLY":     37,
	"SEP":           38,
	"COLON":         39,
	"ANTI_COLON":    40,
	"ASSIGN":        41,
	"UNDETERMINED":  42,
	"COMMENT_ONE":   43,
	"L_BOOL":        44,
	"L_INT":         45,
	"L_STRING":      46,
	"CO_EQUALS":     47,
	"CO_NOT_EQUALS": 48,
	"CO_GT":         49,
	"CO_GT_EQUALS":  50,
	"CO_LT":         51,
	"CO_LT_EQUALS":  52,
	"LITERAL":       53,
	"COMPARISON":    54,
}
