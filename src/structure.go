package main

import "strings"

type Structure struct {
	code     int
	text     string
	line     int
	children []Structure
}

func createStructure(code, text string, line int) Structure {
	return Structure{
		structureCode[code],
		text,
		line,
		[]Structure{},
	}
}

func (st Structure) stringify() string {
	text := ""
	for i := 0; i < len(st.children); i++ {
		text += "\n" + st.children[i].stringify()
	}
	return st.text + strings.ReplaceAll(text, "\n", "\n\t")
}

var structureCode map[string]int = map[string]int{
	// Not implemented
	"ILLEGAL": -1,

	// Statements
	"ST_IMPORT":       0,
	"IF_ELSE_BLOCK":   1,
	"ST_IF":           2,
	"ST_ELIF":         3,
	"ST_ELSE":         4,
	"ST_FOR":          5,
	"ST_WHILE":        6,
	"ST_DECLARATION":  7,
	"ST_MANIPULATION": 8,
	"ST_CALL":         9,
	"ST_FUNCTION":     10,
	"ST_RETURN":       11,

	// Other
	"BLOCK":        32,
	"EXPRESSION":   33,
	"COMPARISON":   34,
	"PROGRAM":      35,
	"ASTERISK":     36,
	"IDENTIFIER":   37,
	"NEWLINE":      38,
	"INDENT":       39,
	"L_PAREN":      40,
	"R_PAREN":      41,
	"L_BLOCK":      42,
	"R_BLOCK":      43,
	"L_SQUIRLY":    44,
	"R_SQUIRLY":    45,
	"SEP":          46,
	"COLON":        47,
	"ANTI_COLON":   48,
	"ASSIGN":       49,
	"UNDETERMINED": 50,
	"COMMENT_ONE":  51,
	"ACCESSOR":     52,
	"FUNC_NAME":    53,
	"ARROW":        54,

	// Keywords
	"K_IMPORT": 64,
	"K_FROM":   65,
	"K_FOR":    66,
	"K_IN":     67,
	"K_IF":     68,
	"K_ELIF":   69,
	"K_ELSE":   70,
	"K_CLASS":  71,
	"K_WHILE":  72,
	"K_DEF":    73,
	"K_RETURN": 74,

	// In-built functions
	"IB_PRINT": 96,
	"IB_RANGE": 97,

	// Bool operands
	"BO_NOT": 128,
	"BO_AND": 129,
	"BO_OR":  130,

	// Math operands
	"MO_PLUS":   160,
	"MO_SUB":    161,
	"MO_MUL":    162,
	"MO_DIV":    163,
	"MO_MODULO": 164,

	// Literal
	"L_BOOL":   192,
	"L_INT":    193,
	"L_STRING": 194,
	"L_NULL":   195,

	// Comparison operands
	"CO_EQUALS":     224,
	"CO_NOT_EQUALS": 225,
	"CO_GT":         226,
	"CO_GT_EQUALS":  227,
	"CO_LT":         228,
	"CO_LT_EQUALS":  229,
}
