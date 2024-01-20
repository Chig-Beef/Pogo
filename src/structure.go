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
	"ST_DECLARATION":  6,
	"ST_MANIPULATION": 7,
	"ST_CALL":         8,

	// Other
	"BLOCK":        32,
	"CALL":         33,
	"EXPRESSION":   34,
	"COMPARISON":   35,
	"PROGRAM":      36,
	"ASTERISK":     37,
	"IDENTIFIER":   38,
	"NEWLINE":      39,
	"INDENT":       40,
	"L_PAREN":      41,
	"R_PAREN":      42,
	"L_BLOCK":      43,
	"R_BLOCK":      44,
	"L_SQUIRLY":    45,
	"R_SQUIRLY":    46,
	"SEP":          47,
	"COLON":        48,
	"ANTI_COLON":   49,
	"ASSIGN":       50,
	"UNDETERMINED": 51,
	"COMMENT_ONE":  52,

	// Keywords
	"K_IMPORT": 64,
	"K_FROM":   65,
	"K_FOR":    66,
	"K_IN":     67,
	"K_IF":     68,
	"K_ELIF":   69,
	"K_ELSE":   70,

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

	// Comparison operands
	"CO_EQUALS":     224,
	"CO_NOT_EQUALS": 225,
	"CO_GT":         226,
	"CO_GT_EQUALS":  227,
	"CO_LT":         228,
	"CO_LT_EQUALS":  229,
}
