package main

import (
	"errors"
	"log"
)

type Parser struct {
	curPos   int
	curToken Token
	source   []Token
	//	markers  []int
}

/*
func (p *Parser) setMarker() {
	p.markers = append(p.markers, p.curPos)
}

func (p *Parser) gotoMarker() {
	p.curPos = p.markers[len(p.markers)-1]
	p.curToken = p.source[p.curPos]
	p.markers = p.markers[:len(p.markers)-1]
}
*/

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

	for p.curToken.code == tokenCode["COMMENT_ONE"] || p.curToken.code == tokenCode["NEWLINE"] {
		sc := structureCode["COMMENT_ONE"]
		if p.curToken.code == tokenCode["NEWLINE"] {
			sc = structureCode["NEWLINE"]
		}
		sts = append(sts, Structure{[]Structure{}, sc, p.curToken.text})
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
					output = append(output, Token{tokenCode["ANTI_COLON"], ":"})
				}
			}
		}

		if input[i].code == tokenCode["INDENT"] {
			continue
		}

		output = append(output, input[i])
	}

	c_count := 0
	ac_count := 0
	for i := 0; i < len(output); i++ {
		if output[i].code == tokenCode["COLON"] {
			c_count++
		} else if output[i].code == tokenCode["ANTI_COLON"] {
			ac_count++
		}
	}

	for i := 0; i < c_count-ac_count-1; i++ {
		output = append(output, Token{tokenCode["ANTI_COLON"], ":"})
	}

	return append(output, Token{tokenCode["NEWLINE"], "NEWLINE"})
}

func (p *Parser) parse(input []Token) Structure {
	if len(input) == 0 {
		log.Fatal("[Parse (parse)] Missing input")
	}

	p.source = input
	p.curPos = 0
	p.curToken = p.source[p.curPos]

	s, err := p.program()
	if err != nil {
		log.Fatal(err.Error())
	}

	return s
}

func (p *Parser) program() (Structure, error) {
	program := Structure{[]Structure{}, structureCode["PROGRAM"], "PROGRAM"}

	for p.curPos < len(p.source) {
		statement, err := p.statement()
		if err != nil {
			return program, err
		}
		program.children = append(program.children, statement)

		program.children = append(program.children, p.nextTokenNoNotes()...)

	}

	return program, nil
}

func (p *Parser) statement() (Structure, error) {
	var s Structure

	if p.curToken.code == tokenCode["K_IMPORT"] {
		s = Structure{[]Structure{}, structureCode["ST_IMPORT"], ""}
		s.children = append(s.children, Structure{[]Structure{}, structureCode["K_IMPORT"], p.curToken.text})
		p.nextToken()

		temp, err := p.checkToken("IDENTIFIER")
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
	} else if p.curToken.code == tokenCode["K_FROM"] {
		s = Structure{[]Structure{}, structureCode["ST_IMPORT"], "ST_IMPORT"}
		s.children = append(s.children, Structure{[]Structure{}, structureCode["K_FROM"], p.curToken.text})
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
				return s, errors.New("[Parse (parse-ST_IMPORT)] Expected ASTERISK, got " + p.curToken.text)
			}
			temp.code = structureCode["ASTERISK"]
		}
		s.children = append(s.children, temp)

		if p.curToken.code != tokenCode["UNDETERMINED"] {
			if p.curToken.code != tokenCode["IDENTIFIER"] {
				return s, errors.New("[Parse (parse-ST_IMPORT)] Expected ASTERISK, got " + p.curToken.text)
			}
			s.children = append(s.children, Structure{[]Structure{}, structureCode["IDENTIFIER"], p.curToken.text})
		} else {
			if p.curToken.text == "*" {
				s.children = append(s.children, Structure{[]Structure{}, structureCode["ASTERISK"], p.curToken.text})
			} else {
				return s, errors.New("[Parse (parse-ST_IMPORT)] Expected ASTERISK, got " + p.curToken.text)
			}
		}
	} else if p.curToken.code == tokenCode["COMMENT_ONE"] {
		s = Structure{[]Structure{}, structureCode["COMMENT_ONE"], p.curToken.text}
	} else if p.curToken.code == tokenCode["K_FOR"] {
		s = Structure{[]Structure{}, structureCode["ST_FOR"], "ST_FOR"}

		s.children = append(s.children, Structure{[]Structure{}, structureCode["K_FOR"], p.curToken.text})
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
		s.children = append(s.children, temps...)

		temp, err = p.block()
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
	} else if p.curToken.code == tokenCode["K_IF"] {
		s = Structure{[]Structure{}, structureCode["ST_IF_ELSE_BLOCK"], "ST_IF_ELSE_BLOCK"}

		temp, err := p.s_if()
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
		p.nextToken()

		for p.curToken.code == tokenCode["K_ELIF"] {
			temp, err = p.s_elif()
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temp)
			p.nextToken()
		}

		temp, err = p.s_else()
		if err != nil {
			return s, err
		}
		s.children = append(s.children, temp)
	} else if p.curToken.code == tokenCode["IB_PRINT"] {
		s = Structure{[]Structure{}, structureCode["ST_CALL"], "ST_CALL"}
		s.children = append(s.children, Structure{[]Structure{}, structureCode["IB_PRINT"], p.curToken.text})
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
			s = Structure{[]Structure{}, structureCode["ST_DECLARATION"], "ST_DECLARATION"}
			s.children = append(s.children, Structure{[]Structure{}, structureCode["IDENTIFIER"], p.curToken.text})
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

			temp, err := p.checkTokenChoices([]string{
				"L_BOOL",
				"L_INT",
				"L_STRING",
			})
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temp)
		} else if p.peek().code == tokenCode["ASSIGN"] {
			s = Structure{[]Structure{}, structureCode["ST_MANIPULATION"], "ST_MANIPULATION"}
			s.children = append(s.children, Structure{[]Structure{}, structureCode["IDENTIFIER"], p.curToken.text})
			p.nextToken()

			temp, err := p.checkToken("ASSIGN")
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temp)
			p.nextToken()

			temp, err = p.checkTokenChoices([]string{
				"L_BOOL",
				"L_INT",
				"L_STRING",
			})
			if err != nil {
				return s, err
			}
			s.children = append(s.children, temp)
		}
	}

	if len(s.children) == 0 {
		s = Structure{[]Structure{}, structureCode["ILLEGAL"], "ILLEGAL"}
	}

	return s, nil
}

func (p *Parser) block() (Structure, error) {
	block := Structure{[]Structure{}, structureCode["BLOCK"], "BLOCK"}

	for p.curPos < len(p.source) {
		statement, err := p.statement()
		if err != nil {
			return block, err
		}
		block.children = append(block.children, statement)

		p.nextToken()

		if p.curToken.code == tokenCode["ANTI_COLON"] {
			break
		}

		block.children = append(block.children, p.nextTokenNoNotes()...)

	}

	block.children = append(block.children, Structure{[]Structure{}, structureCode["ANTI_COLON"], ":"})

	return block, nil
}

func (p *Parser) s_if() (Structure, error) {
	s := Structure{[]Structure{}, structureCode["ST_IF"], "ST_IF"}

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
	s.children = append(s.children, temps...)

	temp, err = p.block()
	if err != nil {
		return s, err
	}
	s.children = append(s.children, temp)

	return s, nil
}

func (p *Parser) s_elif() (Structure, error) {
	return Structure{}, nil
}

func (p *Parser) s_else() (Structure, error) {
	return Structure{}, nil
}

func (p *Parser) expression() (Structure, error) {
	s := Structure{[]Structure{}, structureCode["EXPRESSION"], "EXPRESSION"}

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

	return s, nil
}

func (p *Parser) comparison() (Structure, error) {
	s := Structure{[]Structure{}, structureCode["COMPARISON"], "COMPARISON"}

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
		return s, err
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
	structures := []Structure{}
	for i := 0; i < len(tokenKeys); i++ {
		temp, err := p.checkToken(tokenKeys[i])
		if err != nil {
			return structures, err
		}
		structures = append(structures, temp)
		p.nextToken()
	}
	return structures, nil
}

func (p *Parser) checkTokenChoices(tokenKeys []string) (Structure, error) {
	for i := 0; i < len(tokenKeys); i++ {
		if p.curToken.code == tokenCode[tokenKeys[i]] {
			return Structure{[]Structure{}, structureCode[tokenKeys[i]], p.curToken.text}, nil
		}
	}
	errText := ""
	for i := 0; i < len(tokenKeys); i++ {
		errText += tokenKeys[i]
		errText += " or "
	}
	errText = errText[:len(errText)-4]
	return Structure{}, errors.New("[Parse (checkToken)] Expected " + errText + ", got " + p.curToken.text)
}

func (p *Parser) checkToken(tokenKey string) (Structure, error) {
	if p.curToken.code == tokenCode[tokenKey] {
		return Structure{[]Structure{}, structureCode[tokenKey], p.curToken.text}, nil
	}
	return Structure{}, errors.New("[Parse (checkToken)] Expected " + tokenKey + ", got " + p.curToken.text)
}
