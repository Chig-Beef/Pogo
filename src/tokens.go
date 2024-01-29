package main

type Token struct {
	code int
	text string
	line int
}

var tokenCode map[string]int = map[string]int{
	// Not implemented
	"ILLEGAL": -1,

	// Keywords
	"K_IMPORT": 0,
	"K_FROM":   1,
	"K_FOR":    2,
	"K_IN":     3,
	"K_IF":     4,
	"K_ELIF":   5,
	"K_ELSE":   6,
	"K_CLASS":  7,
	"K_WHILE":  8,
	"K_DEF":    9,
	"K_RETURN": 10,

	// In-Built Funcs
	"IB_PRINT": 32,
	"IB_RANGE": 33,

	// Bool operands
	"BO_NOT": 64,
	"BO_AND": 65,
	"BO_OR":  66,

	// Math operands
	"MO_PLUS":   66,
	"MO_SUB":    67,
	"MO_MUL":    68,
	"MO_DIV":    69,
	"MO_MODULO": 70,

	// Other
	"IDENTIFIER":   128,
	"NEWLINE":      129,
	"INDENT":       130,
	"L_PAREN":      131,
	"R_PAREN":      132,
	"L_BLOCK":      133,
	"R_BLOCK":      134,
	"L_SQUIRLY":    135,
	"R_SQUIRLY":    136,
	"SEP":          137,
	"COLON":        138,
	"ANTI_COLON":   139,
	"ASSIGN":       140,
	"UNDETERMINED": 141,
	"COMMENT_ONE":  142,
	"ACCESSOR":     143,
	"FUNC_NAME":    144,
	"ARROW":        145,

	// Literals
	"L_BOOL":   160,
	"L_INT":    161,
	"L_STRING": 162,
	"L_NULL":   163,

	// Comparison Operands
	"CO_EQUALS":     192,
	"CO_NOT_EQUALS": 193,
	"CO_GT":         194,
	"CO_GT_EQUALS":  195,
	"CO_LT":         196,
	"CO_LT_EQUALS":  197,
}
