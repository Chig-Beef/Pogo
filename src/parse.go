package main

import (
	"errors"
	"log"
)

type Parser struct {
	curPos    int
	curToken  Token
	source    []Token
	markers   []int
	functions []Structure
	funcLine  []string
}

func (p *Parser) setMarker() {
	p.markers = append(p.markers, p.curPos)
}

func (p *Parser) gotoMarker() {
	p.curPos = p.markers[len(p.markers)-1]
	p.curToken = p.source[p.curPos]
	p.markers = p.markers[:len(p.markers)-1]
}

func (p *Parser) nextToken() {
	p.curPos++
	if p.curPos >= len(p.source) {
		p.curToken = Token{} // Nil
	} else {
		p.curToken = p.source[p.curPos]
	}
}

func (p *Parser) rollBack() {
	p.curPos--
	if p.curPos < len(p.source) {
		p.curToken = Token{} // Nil
	} else {
		p.curToken = p.source[p.curPos]
	}
}

func (p *Parser) nextTokenNoNotes() []Structure {
	p.nextToken()
	sts := []Structure{}

	for p.curToken.code == tokenCode["COMMENT_ONE"] || p.curToken.code == tokenCode["COMMENT_MULTI"] || p.curToken.code == tokenCode["NEWLINE"] {
		sc := structureCode["COMMENT_ONE"]
		if p.curToken.code == tokenCode["NEWLINE"] {
			sc = structureCode["NEWLINE"]
		} else if p.curToken.code == tokenCode["COMMENT_MULTI"] {
			sc = structureCode["COMMENT_MULTI"]
		}
		sts = append(sts, Structure{sc, p.curToken.text, p.curToken.line, []Structure{}})
		p.nextToken()
	}
	return sts
}

func (p *Parser) peek() Token {
	if p.curPos >= len(p.source)-1 {
		return Token{}
	}
	return p.source[p.curPos+1]
}

func (p *Parser) replaceIndents(input []Token) []Token {
	indents := []int{0}
	curIndex := 0

	for i := 0; i < len(input); i++ {
		if input[i].code == tokenCode["NEWLINE"] {
			curIndex++
			indents = append(indents, 0)
		} else if input[i].code == tokenCode["INDENT"] {
			indents[curIndex]++
		}
	}

	output := []Token{}
	curIndex = 0
	for i := 0; i < len(input); i++ {
		if input[i].code == tokenCode["NEWLINE"] {
			curIndex++
			if indents[curIndex] < indents[curIndex-1] {
				for j := 0; j < indents[curIndex-1]-indents[curIndex]; j++ {
					output = append(output, Token{tokenCode["ANTI_COLON"], ":", input[i].line})
				}
			}
		}

		if input[i].code == tokenCode["INDENT"] {
			continue
		}

		output = append(output, input[i])
	}

	for i := 0; i < indents[len(indents)-1]; i++ {
		output = append(output, Token{tokenCode["ANTI_COLON"], ":", len(indents) - 1})
	}

	return append(output, Token{tokenCode["NEWLINE"], "NEWLINE", len(indents) - 1})
}

func (p *Parser) checkImport(program Structure) (Structure, error) {
	if program.children[0].code != structureCode["ST_IMPORT"] {
		return program, errors.New("(Parse [checkImport]) Source should start with \"from GoType import *\"")
	}
	if program.children[0].children[0].code != structureCode["K_FROM"] {
		return program, errors.New("(Parse [checkImport]) Source should start with \"from GoType import *\"")
	}
	if program.children[0].children[1].text != "GoType" {
		return program, errors.New("(Parse [checkImport]) Source should start with \"from GoType import *\"")
	}
	program.children = program.children[1:]
	return program, nil
}

func (p *Parser) parse(input []Token) Structure {
	p.funcLine = []string{"parse.go", "parse"}

	if len(input) == 0 {
		log.Fatal(createError(p.funcLine, "Missing input", 0))
	}

	p.source = input
	p.curPos = 0
	p.curToken = p.source[p.curPos]

	s, err := p.program()
	if err != nil {
		log.Fatal(err.Error())
	}

	s, err = p.checkImport(s)
	if err != nil {
		log.Fatal(err.Error())
	}

	return s
}

func (p *Parser) program() (Structure, error) {
	p.funcLine = append(p.funcLine, "program")
	program := createStructure("PROGRAM", "PROGRAM", 0)

	for p.curPos < len(p.source) {
		statement, err := p.statement()
		if err != nil {
			return program, err
		}
		program.children = append(program.children, statement)

		program.children = append(program.children, p.nextTokenNoNotes()...)

	}

	p.funcLine = p.funcLine[:len(p.funcLine)-1]

	return program, nil
}

func (p *Parser) statement() (Structure, error) {
	p.funcLine = append(p.funcLine, "statement")
	var s Structure

	if p.curToken.code == tokenCode["K_IMPORT"] {
		s = createStructure("ST_IMPORT", "", p.curToken.line)
		s.children = append(s.children, createStructure("K_IMPORT", p.curToken.text, p.curToken.line))
		p.nextToken()

		temp, err := p.checkToken("IDENTIFIER")
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
	} else if p.curToken.code == tokenCode["K_FROM"] {
		s = createStructure("ST_IMPORT", "ST_IMPORT", p.curToken.line)
		s.children = append(s.children, createStructure("K_FROM", p.curToken.text, p.curToken.line))
		p.nextToken()

		temps, err := p.checkTokenRange([]string{
			"IDENTIFIER",
			"K_IMPORT",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temps...)

		temp, err := p.checkToken("UNDETERMINED")
		if err != nil {
			temp, err = p.checkToken("IDENTIFIER")
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temp)
		} else {
			if p.curToken.text != "*" {
				return s, createError(p.funcLine, "Expected ASTERISK, got "+p.curToken.text, p.curToken.line)
			}
			temp.code = structureCode["ASTERISK"]
		}
		s.children = append(s.children, temp)

		if p.curToken.code != tokenCode["UNDETERMINED"] {
			if p.curToken.code != tokenCode["IDENTIFIER"] {
				return s, createError(p.funcLine, "Expected ASTERISK, got "+p.curToken.text, p.curToken.line)
			}
			s.children = append(s.children, createStructure("IDENTIFIER", p.curToken.text, p.curToken.line))
		} else {
			if p.curToken.text == "*" {
				s.children = append(s.children, createStructure("ASTERISK", p.curToken.text, p.curToken.line))
			} else {
				return s, createError(p.funcLine, "Expected ASTERISK, got"+p.curToken.text, p.curToken.line)
			}
		}
	} else if p.curToken.code == tokenCode["COMMENT_ONE"] {
		s = createStructure("COMMENT_ONE", p.curToken.text, p.curToken.line)
	} else if p.curToken.code == tokenCode["COMMENT_MULTI"] {
		s = createStructure("COMMENT_MULTI", p.curToken.text, p.curToken.line)
		s.children = append(s.children, createStructure("NEWLINE", "NEWLINE", p.curToken.line))
	} else if p.curToken.code == tokenCode["K_FOR"] {
		s = createStructure("ST_FOR", "ST_FOR", p.curToken.line)

		s.children = append(s.children, createStructure("K_FOR", p.curToken.text, p.curToken.line))
		p.nextToken()

		temps, err := p.checkTokenRange([]string{
			"IDENTIFIER",
			"K_IN",
			"IB_RANGE",
			"L_PAREN",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temps...)

		temp, err := p.checkTokenChoices([]string{
			"L_INT",
			"IDENTIFIER",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
		p.nextToken()

		temp, err = p.checkToken("SEP")
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
		p.nextToken()

		temp, err = p.checkTokenChoices([]string{
			"L_INT",
			"IDENTIFIER",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
		p.nextToken()

		temps, err = p.checkTokenRange([]string{
			"R_PAREN",
			"COLON",
			"NEWLINE",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temps[:2]...)

		temp, err = p.block()
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
	} else if p.curToken.code == tokenCode["K_WHILE"] {
		s = createStructure("ST_WHILE", "ST_WHILE", p.curToken.line)

		s.children = append(s.children, createStructure("K_WHILE", p.curToken.text, p.curToken.line))
		p.nextToken()

		temp, err := p.comparison()
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
		p.nextToken()

		temps, err := p.checkTokenRange([]string{
			"COLON",
			"NEWLINE",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temps...)

		temp, err = p.block()
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
	} else if p.curToken.code == tokenCode["K_IF"] {
		s = createStructure("ST_IF_ELSE_BLOCK", "ST_IF_ELSE_BLOCK", p.curToken.line)

		temp, err := p.s_if()
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)

		p.setMarker()
		p.nextTokenNoNotes()

		extra_found := false
		for p.curToken.code == tokenCode["K_ELIF"] {
			extra_found = true
			temp, err = p.s_elif()
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temp)
			p.nextTokenNoNotes()
		}

		if p.curToken.code == tokenCode["K_ELSE"] {
			extra_found = true
			temp, err = p.s_else()
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temp)
		}
		if !extra_found {
			p.gotoMarker()
		}
	} else if p.curToken.code == tokenCode["IB_PRINT"] {
		s = createStructure("ST_CALL", "ST_CALL", p.curToken.line)
		s.children = append(s.children, createStructure("IB_PRINT", p.curToken.text, p.curToken.line))
		p.nextToken()

		temp, err := p.checkToken("L_PAREN")
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
		p.nextToken()

		temp, err = p.checkTokenChoices([]string{
			"L_BOOL",
			"L_INT",
			"L_STRING",
			"IDENTIFIER",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
		p.nextToken()

		temp, err = p.checkToken("R_PAREN")
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
	} else if p.curToken.code == tokenCode["IDENTIFIER"] {
		if p.peek().code == tokenCode["COLON"] {
			s = createStructure("ST_DECLARATION", "ST_DECLARATION", p.curToken.line)
			s.children = append(s.children, createStructure("IDENTIFIER", p.curToken.text, p.curToken.line))
			p.nextToken()

			temps, err := p.checkTokenRange([]string{
				"COLON",
				"IDENTIFIER",
				"ASSIGN",
			})
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temps...)

			p.setMarker()
			temp, err := p.call()
			if err != nil {
				p.gotoMarker()
				temp, err = p.checkTokenChoices([]string{
					"L_BOOL",
					"L_INT",
					"L_STRING",
				})
				if err != nil {
					return s, err
				}
			}
			s.children = append(s.children, temp)
		} else if p.peek().code == tokenCode["ASSIGN"] {
			s = createStructure("ST_MANIPULATION", "ST_MANIPULATION", p.curToken.line)
			s.children = append(s.children, createStructure("IDENTIFIER", p.curToken.text, p.curToken.line))
			p.nextToken()

			temp, err := p.checkToken("ASSIGN")
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temp)
			p.nextToken()

			temp, err = p.expression()
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temp)
		} else if p.peek().code == tokenCode["L_PAREN"] {
			temp, err := p.call()
			if err != nil {
				return s, err
			}
			s = temp
		}
	} else if p.curToken.code == tokenCode["K_DEF"] {
		s = createStructure("ST_FUNCTION", "ST_FUNCTION", p.curToken.line)
		s.children = append(s.children, createStructure("K_DEF", p.curToken.text, p.curToken.line))
		p.nextToken()

		temp, err := p.checkToken("IDENTIFIER")
		if err != nil {
			return s, err
		}
		temp.code = structureCode["FUNC_NAME"]
		s.children = append(s.children, temp)
		p.nextToken()

		temp, err = p.checkToken("L_PAREN")
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)

		p.nextToken()

		for p.curToken.code == tokenCode["IDENTIFIER"] {
			temps, err := p.checkTokenRange([]string{
				"IDENTIFIER",
				"COLON",
				"IDENTIFIER",
			})
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temps...)

			temp, err = p.checkToken("SEP")
			if err != nil {
				break
			}
			s.children = append(s.children, temp)
			p.nextToken()
		}

		temps, err := p.checkTokenRange([]string{
			"R_PAREN",
			"ARROW",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temps...)
		p.nextToken()

		temp, err = p.checkTokenChoices([]string{
			"IDENTIFIER",
			"L_NULL",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
		p.nextToken()

		temps, err = p.checkTokenRange([]string{
			"COLON",
			"NEWLINE",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temps[0])

		temp, err = p.block()
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)

		p.functions = append(p.functions, s)

		p.funcLine = p.funcLine[:len(p.funcLine)-1]
		return createStructure("NEWLINE", "NEWLINE", p.curToken.line), nil
	} else if p.curToken.code == tokenCode["K_RETURN"] {
		s = createStructure("ST_RETURN", "ST_RETURN", p.curToken.line)
		s.children = append(s.children, createStructure("K_RETURN", p.curToken.text, p.curToken.line))
		p.nextToken()

		temp, err := p.checkTokenChoices([]string{
			"IDENTIFIER",
			"L_BOOL",
			"L_INT",
			"L_STRING",
		})
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
	}

	if len(s.children) == 0 {
		s = createStructure("ILLEGAL", "ILLEGAL + "+p.curToken.text, p.curToken.line)
	}

	p.funcLine = p.funcLine[:len(p.funcLine)-1]
	return s, nil
}

func (p *Parser) block() (Structure, error) {
	p.funcLine = append(p.funcLine, "block")
	block := createStructure("BLOCK", "BLOCK", p.curToken.line)

	for p.curPos < len(p.source) {
		statement, err := p.statement()
		if err != nil {
			return block, err
		}
		block.children = append(block.children, statement)

		//p.nextToken()

		if p.peek().code == tokenCode["ANTI_COLON"] {
			p.nextToken()
			break
		}

		block.children = append(block.children, p.nextTokenNoNotes()...)

	}

	block.children = append(block.children, createStructure("ANTI_COLON", ":", p.curToken.line))

	p.funcLine = p.funcLine[:len(p.funcLine)-1]
	return block, nil
}

func (p *Parser) call() (Structure, error) {
	var err error

	s := createStructure("ST_CALL", "ST_CALL", p.curToken.line)

	temp, err := p.checkToken("IDENTIFIER")
	if err != nil {
		return s, err
	}
	temp.code = structureCode["FUNC_NAME"]
	s.children = append(s.children, temp)
	p.nextToken()

	temp, err = p.checkToken("L_PAREN")
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temp, err = p.checkTokenChoices([]string{
		"IDENTIFIER",
		"L_BOOL",
		"L_INT",
		"L_STRING",
	})

	for err == nil {
		if p.peek().code == tokenCode["L_PAREN"] { // We're dealing with a call
			temp, err = p.call()
			if err != nil {
				return s, err
			}
		}

		s.children = append(s.children, temp)
		p.nextToken()

		temp, err = p.checkToken("SEP")
		if err != nil {
			break
		}
		s.children = append(s.children, temp)
		p.nextToken()

		temp, err = p.checkTokenChoices([]string{
			"IDENTIFIER",
			"L_BOOL",
			"L_INT",
			"L_STRING",
		})
	}

	temp, err = p.checkToken("R_PAREN")
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)

	return s, nil
}

func (p *Parser) s_if() (Structure, error) {
	p.funcLine = append(p.funcLine, "s_if")
	s := createStructure("ST_IF", "ST_IF", p.curToken.line)

	temp, err := p.checkToken("K_IF")
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temp, err = p.comparison()
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temps, err := p.checkTokenRange([]string{
		"COLON",
		"NEWLINE",
	})
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temps[0])

	temp, err = p.block()
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)

	p.funcLine = p.funcLine[:len(p.funcLine)-1]
	return s, nil
}

func (p *Parser) s_elif() (Structure, error) {
	p.funcLine = append(p.funcLine, "s_elif")
	s := createStructure("ST_ELIF", "ST_ELIF", p.curToken.line)

	temp, err := p.checkToken("K_ELIF")
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temp, err = p.comparison()
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temps, err := p.checkTokenRange([]string{
		"COLON",
		"NEWLINE",
	})
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temps[0])

	temp, err = p.block()
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)

	p.funcLine = p.funcLine[:len(p.funcLine)-1]
	return s, nil
}

func (p *Parser) s_else() (Structure, error) {
	p.funcLine = append(p.funcLine, "s_else")
	s := createStructure("ST_ELSE", "ST_ELSE", p.curToken.line)

	temp, err := p.checkToken("K_ELSE")
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temps, err := p.checkTokenRange([]string{
		"COLON",
		"NEWLINE",
	})
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temps[0])

	temp, err = p.block()
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)

	p.funcLine = p.funcLine[:len(p.funcLine)-1]
	return s, nil
}

func (p *Parser) expression() (Structure, error) {
	p.funcLine = append(p.funcLine, "expression")
	s := createStructure("EXPRESSION", "EXPRESSION", p.curToken.line)

	temp, err := p.checkTokenChoices([]string{
		"L_BOOL",
		"L_INT",
		"L_STRING",
		"IDENTIFIER",
	})
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temp, err = p.checkTokenChoices([]string{
		"MO_PLUS",
		"MO_SUB",
		"MO_MUL",
		"MO_DIV",
		"MO_MODULO",
	})
	if err != nil {
		p.rollBack()
		p.funcLine = p.funcLine[:len(p.funcLine)-1]
		return s, nil // Could be a single literal, so we don't error
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temp, err = p.checkTokenChoices([]string{
		"L_BOOL",
		"L_INT",
		"L_STRING",
		"IDENTIFIER",
	})
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)

	p.funcLine = p.funcLine[:len(p.funcLine)-1]
	return s, nil
}

func (p *Parser) comparison() (Structure, error) {
	p.funcLine = append(p.funcLine, "comparison")
	s := createStructure("COMPARISON", "COMPARISON", p.curToken.line)

	temp, err := p.expression()
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temp, err = p.checkTokenChoices([]string{
		"CO_EQUALS",
		"CO_NOT_EQUALS",
		"CO_GT",
		"CO_GT_EQUALS",
		"CO_LT",
		"CO_LT_EQUALS",
	})
	if err != nil {
		p.rollBack()
		p.funcLine = p.funcLine[:len(p.funcLine)-1]
		return s, nil // Could be a single literal, so we don't error
	}
	s.children = append(s.children, temp)
	p.nextToken()

	temp, err = p.expression()
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)

	return s, nil
}

func (p *Parser) checkTokenRange(tokenKeys []string) ([]Structure, error) {
	p.funcLine = append(p.funcLine, "checkTokenRange")
	structures := []Structure{}
	for i := 0; i < len(tokenKeys); i++ {
		temp, err := p.checkToken(tokenKeys[i])
		if err != nil {
			return structures, err
		}
		structures = append(structures, temp)
		p.nextToken()
	}
	p.funcLine = p.funcLine[:len(p.funcLine)-1]
	return structures, nil
}

func (p *Parser) checkTokenChoices(tokenKeys []string) (Structure, error) {
	p.funcLine = append(p.funcLine, "checkTokenChoices")
	for i := 0; i < len(tokenKeys); i++ {
		if p.curToken.code == tokenCode[tokenKeys[i]] {
			p.funcLine = p.funcLine[:len(p.funcLine)-1]
			return createStructure(tokenKeys[i], p.curToken.text, p.curToken.line), nil
		}
	}
	errText := ""
	for i := 0; i < len(tokenKeys); i++ {
		errText += tokenKeys[i]
		errText += " or "
	}
	errText = errText[:len(errText)-4]
	return Structure{}, createError(p.funcLine, "Expected "+errText+", got "+p.curToken.text, p.curToken.line)
}

func (p *Parser) checkToken(tokenKey string) (Structure, error) {
	p.funcLine = append(p.funcLine, "checkToken")
	if p.curToken.code == tokenCode[tokenKey] {
		p.funcLine = p.funcLine[:len(p.funcLine)-1]
		return createStructure(tokenKey, p.curToken.text, p.curToken.line), nil
	}
	return Structure{}, createError(p.funcLine, "Expected "+tokenKey+", got "+p.curToken.text, p.curToken.line)
}
