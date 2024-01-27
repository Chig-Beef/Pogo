package main

import (
	"log"
	"strconv"
	"unicode"
)

type Lexer struct {
	curPos  int
	curChar byte
	source  []byte
	line    int
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
	l.line = 1

	tokens := []Token{}

	for l.curPos < len(l.source) {
		var token Token

		// Math Operands
		if l.curChar == '+' {
			token = Token{tokenCode["MO_PLUS"], "+", l.line}
		} else if l.curChar == '-' {
			token = Token{tokenCode["MO_SUN"], "-", l.line}
		} else if l.curChar == '*' {
			token = Token{tokenCode["UNDETERMINED"], "*", l.line} // Could be for import
		} else if l.curChar == '/' {
			token = Token{tokenCode["MO_DIV"], "/", l.line}
		} else if l.curChar == '%' {
			token = Token{tokenCode["MO_MODULO"], "%", l.line}
		}

		// Parens
		if l.curChar == '(' {
			token = Token{tokenCode["L_PAREN"], "(", l.line}
		} else if l.curChar == ')' {
			token = Token{tokenCode["R_PAREN"], ")", l.line}
		} else if l.curChar == '[' {
			token = Token{tokenCode["L_BLOCK"], ")", l.line}
		} else if l.curChar == ']' {
			token = Token{tokenCode["R_BLOCK"], ")", l.line}
		} else if l.curChar == '{' {
			token = Token{tokenCode["L_SQUIRLY"], ")", l.line}
		} else if l.curChar == '}' {
			token = Token{tokenCode["R_SQUIRLY"], ")", l.line}
		}

		// Other
		if l.curChar == '\r' && l.peek() == '\n' {
			token = Token{tokenCode["NEWLINE"], "NEWLINE", l.line}
			l.nextChar()
			l.line++
		} else if l.curChar == ',' {
			token = Token{tokenCode["SEP"], ",", l.line}
		} else if l.curChar == ':' {
			token = Token{tokenCode["COLON"], ":", l.line}
		} else if l.curChar == '#' {
			start := l.curPos
			for l.peek() != '\r' && l.peek() != '\n' {
				l.nextChar()
			}
			note := string(l.source[start : l.curPos+1])
			token = Token{tokenCode["COMMENT_ONE"], note, l.line}
		} else if l.curChar == ' ' {
			if string(l.source[l.curPos:l.curPos+4]) == "    " {
				token = Token{tokenCode["INDENT"], "    ", l.line}
				l.nextChar()
				l.nextChar()
				l.nextChar()
			}
		}

		// Comparison Operands
		if l.curChar == '=' {
			if l.peek() == '=' {
				token = Token{tokenCode["CO_EQUALS"], "==", l.line}
				l.nextChar()
			} else {
				token = Token{tokenCode["ASSIGN"], "=", l.line}
			}
		} else if l.curChar == '!' {
			if l.peek() == '=' {
				token = Token{tokenCode["CO_NOT_EQUALS"], "!=", l.line}
				l.nextChar()
			}
		} else if l.curChar == '>' {
			if l.peek() == '=' {
				token = Token{tokenCode["CO_GT_EQUALS"], ">=", l.line}
				l.nextChar()
			} else {
				token = Token{tokenCode["CO_GT"], ">", l.line}
			}
		} else if l.curChar == '<' {
			if l.peek() == '=' {
				token = Token{tokenCode["CO_LT_EQUALS"], "<=", l.line}
				l.nextChar()
			} else {
				token = Token{tokenCode["CO_LT"], "<", l.line}
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
				token = Token{tokenCode["K_IMPORT"], word, l.line}
			} else if word == "from" {
				token = Token{tokenCode["K_FROM"], word, l.line}
			} else if word == "for" {
				token = Token{tokenCode["K_FOR"], word, l.line}
			} else if word == "while" {
				token = Token{tokenCode["K_WHILE"], word, l.line}
			} else if word == "in" {
				token = Token{tokenCode["K_IN"], word, l.line}
			} else if word == "if" {
				token = Token{tokenCode["K_IF"], word, l.line}
			} else if word == "elif" {
				token = Token{tokenCode["K_ELIF"], word, l.line}
			} else if word == "else" {
				token = Token{tokenCode["K_ELSE"], word, l.line}
			} else if word == "def" {
				token = Token{tokenCode["K_DEF"], word, l.line}
			} else if word == "return" {
				token = Token{tokenCode["K_RETURN"], word, l.line}
			}

			// In-Built Funcs
			if word == "print" {
				token = Token{tokenCode["IB_PRINT"], word, l.line}
			} else if word == "range" {
				token = Token{tokenCode["IB_RANGE"], word, l.line}
			}

			// Bool operands
			if word == "not" {
				token = Token{tokenCode["BO_NOT"], word, l.line}
			} else if word == "and" {
				token = Token{tokenCode["BO_AND"], word, l.line}
			} else if word == "or" {
				token = Token{tokenCode["BO_OR"], word, l.line}
			}

			// Boolean literal
			if word == "True" || word == "False" {
				token = Token{tokenCode["L_BOOL"], word, l.line}
			}

			// Identifier
			if token == (Token{}) {
				token = Token{tokenCode["IDENTIFIER"], word, l.line}
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
			token = Token{tokenCode["L_STRING"], num, l.line}
		}

		// Integer literal
		if unicode.IsDigit(rune(l.curChar)) {
			start := l.curPos
			for unicode.IsDigit(rune(l.peek())) || l.peek() == '_' || l.peek() == '.' {
				l.nextChar()
			}
			num := string(l.source[start : l.curPos+1])

			if num[len(num)-1] == '_' || num[len(num)-1] == '.' {
				log.Fatal("[Lex (lex)] Numbers must end with a digit on line " + strconv.Itoa(l.line))
				return tokens
			}

			token = Token{tokenCode["L_INT"], num, l.line}
		}

		// Not Implemented
		if token == (Token{}) {
			token = Token{tokenCode["ILLEGAL"], "", l.line}
		}

		tokens = append(tokens, token)

		l.nextCharNoWhiteSpace()
	}

	return tokens
}
