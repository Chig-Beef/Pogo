package main

import (
	"log"
	"unicode"
)

type Lexer struct {
	curPos  int
	curChar byte
	source  []byte
}

func (l *Lexer) nextChar() {
	l.curPos++
	if l.curPos == len(l.source) {
		l.curChar = 0 // Nil
	} else {
		l.curChar = l.source[l.curPos]
	}
}

func (l *Lexer) nextCharNoWhiteSpace() {
	l.nextChar()
	for l.curChar == ' ' {
		if string(l.source[l.curPos:l.curPos+4]) == "    " {
			break
		}
		l.nextChar()
	}
}

func (l *Lexer) peek() byte {
	if l.curPos >= len(l.source)-1 {
		return 0
	}
	return l.source[l.curPos+1]
}

func (l *Lexer) lex(input []byte) []Token {
	if len(input) == 0 {
		log.Fatal("[Lex (lex)] Missing input")
	}

	l.source = input
	l.curPos = 0
	l.curChar = l.source[l.curPos]

	tokens := []Token{}

	for l.curPos < len(l.source) {
		var token Token

		// Math Operands
		if l.curChar == '+' {
			token = Token{tokenCode["MO_PLUS"], "+"}
		} else if l.curChar == '-' {
			token = Token{tokenCode["MO_SUN"], "-"}
		} else if l.curChar == '*' {
			token = Token{tokenCode["UNDETERMINED"], "*"} // Could be for import
		} else if l.curChar == '/' {
			token = Token{tokenCode["MO_DIV"], "/"}
		} else if l.curChar == '%' {
			token = Token{tokenCode["MO_MODULO"], "%"}
		}

		// Parens
		if l.curChar == '(' {
			token = Token{tokenCode["L_PAREN"], "("}
		} else if l.curChar == ')' {
			token = Token{tokenCode["R_PAREN"], ")"}
		} else if l.curChar == '[' {
			token = Token{tokenCode["L_BLOCK"], ")"}
		} else if l.curChar == ']' {
			token = Token{tokenCode["R_BLOCK"], ")"}
		} else if l.curChar == '{' {
			token = Token{tokenCode["L_SQUIRLY"], ")"}
		} else if l.curChar == '}' {
			token = Token{tokenCode["R_SQUIRLY"], ")"}
		}

		// Other
		if l.curChar == '\r' && l.peek() == '\n' {
			token = Token{tokenCode["NEWLINE"], "NEWLINE"}
			l.nextChar()
		} else if l.curChar == ',' {
			token = Token{tokenCode["SEP"], ","}
		} else if l.curChar == ':' {
			token = Token{tokenCode["COLON"], ":"}
		} else if l.curChar == '#' {
			start := l.curPos
			for l.peek() != '\r' && l.peek() != '\n' {
				l.nextChar()
			}
			note := string(l.source[start : l.curPos+1])
			token = Token{tokenCode["COMMENT_ONE"], note}
		} else if l.curChar == ' ' {
			if string(l.source[l.curPos:l.curPos+4]) == "    " {
				token = Token{tokenCode["INDENT"], "    "}
				l.nextChar()
				l.nextChar()
				l.nextChar()
			}
		}

		// Comparison Operands
		if l.curChar == '=' {
			if l.peek() == '=' {
				token = Token{tokenCode["CO_EQUALS"], "=="}
				l.nextChar()
			} else {
				token = Token{tokenCode["ASSIGN"], "="}
			}
		} else if l.curChar == '!' {
			if l.peek() == '=' {
				token = Token{tokenCode["CO_NOT_EQUALS"], "!="}
				l.nextChar()
			}
		} else if l.curChar == '>' {
			if l.peek() == '=' {
				token = Token{tokenCode["CO_GT_EQUALS"], ">="}
				l.nextChar()
			} else {
				token = Token{tokenCode["CO_GT"], ">"}
			}
		} else if l.curChar == '<' {
			if l.peek() == '=' {
				token = Token{tokenCode["CO_LT_EQUALS"], "<="}
				l.nextChar()
			} else {
				token = Token{tokenCode["CO_LT"], "<"}
			}
		}

		// Words
		if unicode.IsLetter(rune(l.curChar)) {
			start := l.curPos
			for unicode.IsLetter(rune(l.peek())) {
				l.nextChar()
			}
			word := string(l.source[start : l.curPos+1])

			// Keywords
			if word == "import" {
				token = Token{tokenCode["K_IMPORT"], word}
			} else if word == "from" {
				token = Token{tokenCode["K_FROM"], word}
			} else if word == "for" {
				token = Token{tokenCode["K_FOR"], word}
			} else if word == "in" {
				token = Token{tokenCode["K_IN"], word}
			} else if word == "if" {
				token = Token{tokenCode["K_IF"], word}
			} else if word == "elif" {
				token = Token{tokenCode["K_ELIF"], word}
			} else if word == "else" {
				token = Token{tokenCode["K_ELSE"], word}
			}

			// In-Built Funcs
			if word == "print" {
				token = Token{tokenCode["IB_PRINT"], word}
			} else if word == "range" {
				token = Token{tokenCode["IB_RANGE"], word}
			}

			// Bool operands
			if word == "not" {
				token = Token{tokenCode["BO_NOT"], word}
			} else if word == "and" {
				token = Token{tokenCode["BO_AND"], word}
			} else if word == "or" {
				token = Token{tokenCode["BO_OR"], word}
			}

			// Boolean literal
			if word == "True" || word == "False" {
				token = Token{tokenCode["L_BOOL"], word}
			}

			// Identifier
			if token == (Token{}) {
				token = Token{tokenCode["IDENTIFIER"], word}
			}
		}

		// String literal
		if l.curChar == '"' {
			start := l.curPos
			l.nextChar()
			for l.curChar != '"' {
				l.nextChar()
			}
			num := string(l.source[start : l.curPos+1])
			token = Token{tokenCode["L_STRING"], num}
		}

		// Integer literal
		if unicode.IsDigit(rune(l.curChar)) {
			start := l.curPos
			for unicode.IsDigit(rune(l.peek())) || l.peek() == '_' || l.peek() == '.' {
				l.nextChar()
			}
			num := string(l.source[start : l.curPos+1])

			if num[len(num)-1] == '_' || num[len(num)-1] == '.' {
				log.Fatal("[Lex (lex)] Numbers must end with a digit")
				return tokens
			}

			token = Token{tokenCode["L_INT"], num}
		}

		// Not Implemented
		if token == (Token{}) {
			token = Token{tokenCode["ILLEGAL"], ""}
		}

		tokens = append(tokens, token)

		l.nextCharNoWhiteSpace()
	}

	return tokens
}
