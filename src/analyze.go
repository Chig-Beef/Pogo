package main

type Analyzer struct {
}

type Variable struct {
	name    string
	varType string
}

type Function struct {
	name    string
	params  []string
	varType string
}

func (a *Analyzer) analyze(s Structure, vars []Variable, funcs []Function) error {
	//println(s.code)

	if s.code == structureCode["ST_MANIPULATION"] {
		name := s.children[0].text
		var valid bool
		for i := 0; i < len(vars); i++ {
			valid = false
			if name == vars[i].name {
				valid = true
				break
			}
		}
		if !valid {
			return createError([]string{"analyze.go", "analyze:ST_MANIPULATION"}, "An attempt to manipulate an uninitialized variable was made", s.line)
		}
	} else if s.code == structureCode["EXPRESSION"] {
		for i := 0; i < len(s.children); i += 2 {
			if s.children[i].code != structureCode["IDENTIFIER"] {
				continue
			}

			name := s.children[i].text
			var valid bool
			for j := 0; j < len(vars); j++ {
				valid = false
				if name == vars[j].name {
					valid = true
					break
				}
			}
			if !valid {
				return createError([]string{"analyze.go", "analyze:EXPRESSION"}, "An uninitialized variable was used in an expression", s.line)
			}
		}
	} else if s.code == structureCode["COMPARISON"] {
		for i := 0; i < len(s.children); i += 2 {
			if s.children[i].code != structureCode["IDENTIFIER"] {
				continue
			}

			name := s.children[i].text
			var valid bool
			for j := 0; j < len(vars); j++ {
				valid = false
				if name == vars[j].name {
					valid = true
					break
				}
			}
			if !valid {
				return createError([]string{"analyze.go", "analyze:COMPARISON"}, "An uninitialized variable was used in a comparison", s.line)
			}
		}
	} else if s.code == structureCode["ST_CALL"] {
		var fn Function

		name := s.children[0].text
		var valid bool
		for i := 0; i < len(funcs); i++ {
			valid = false
			if name == funcs[i].name {
				valid = true
				fn = funcs[i]
				break
			}
		}
		if !valid {
			return createError([]string{"analyze.go", "analyze:ST_CALL"}, "An attempt to call the non-existent function \""+s.children[0].text+"\" was made", s.line)
		}

		i := 2
		pIndex := 0

		for i+1 < len(s.children) {

			if fn.params[pIndex] == "any" {
				i += 2
				pIndex++
				continue
			}

			if s.children[i].code != structureCode["IDENTIFIER"] {
				switch s.children[i].code {
				case structureCode["L_STRING"]:
					if fn.params[pIndex] != "string" {
						return createError([]string{"analyze.go", "analyze:ST_CALL"}, "Excpected "+fn.params[pIndex]+" got string in function call", s.line)
					}
				case structureCode["L_BOOL"]:
					if fn.params[pIndex] != "bool" {
						return createError([]string{"analyze.go", "analyze:ST_CALL"}, "Excpected "+fn.params[pIndex]+" got bool in function call", s.line)
					}
				case structureCode["L_INT"]:
					valid := false
					types := []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "float32", "float64"}
					for j := 0; j < len(types); j++ {
						if types[j] == fn.params[pIndex] {
							valid = true
							break
						}
					}
					if !valid {
						return createError([]string{"analyze.go", "analyze:ST_CALL"}, "Excpected "+fn.params[pIndex]+" got int in function call", s.line)
					}
				case structureCode["ST_CALL"]:
					var insideCall Function

					name := s.children[i].children[0].text
					var valid bool
					for i := 0; i < len(funcs); i++ {
						valid = false
						if name == funcs[i].name {
							valid = true
							insideCall = funcs[i]
							break
						}
					}
					if !valid {
						return createError([]string{"analyze.go", "analyze:ST_CALL"}, "An attempt to call the non-existent function \""+name+"\" was made", s.line)
					}

					if fn.params[pIndex] != insideCall.varType {
						return createError([]string{"analyze.go", "analyze:ST_CALL"}, "\""+name+"\" is the wrong type, expected "+fn.params[pIndex]+" got "+insideCall.varType, s.line)
					}

				default:
					return createError([]string{"analyze.go", "analyze:ST_CALL"}, "How did you even...? "+fn.name, s.line)
				}
				i += 2
				pIndex++
				continue
			}

			name := s.children[i].text
			var variable Variable
			var valid bool
			for i := 0; i < len(vars); i++ {
				valid = false
				if name == vars[i].name {
					valid = true
					variable = vars[i]
					break
				}
			}
			if !valid {
				return createError([]string{"analyze.go", "analyze:ST_CALL"}, "An uninitialized variable was used in a function call", s.line)
			}

			if variable.varType != fn.params[pIndex] && fn.params[pIndex] != "any" {
				return createError([]string{"analyze.go", "analyze:ST_CALL"}, "Wrong type used in function call", s.line)
			}

			i += 2
			pIndex++
		}
	} else if s.code == structureCode["ST_FUNCTION"] {
		i := 3
		for s.children[i].code == structureCode["IDENTIFIER"] {
			n := s.children[i]   // IDENTIFIER - name
			t := s.children[i+2] // IDENTIFIER - type
			v := Variable{n.text, t.text}
			vars = append(vars, v)

			i += 4 // Next variable, otherwise this will end up on an anti-colon
		}
	}

	if s.code == structureCode["ST_DECLARATION"] {
		name := s.children[0].text
		var variable Variable
		var valid bool
		for i := 0; i < len(vars); i++ {
			valid = false
			if name == vars[i].name {
				valid = true
				variable = vars[i]
				break
			}
		}

		if !valid {
			return createError([]string{"analyze.go", "analyze:ST_DECLARATION"}, "An uninitialized variable was used in a declaration", s.line)
		}

		i := 0
		for i < len(s.children[4].children) {
			if s.children[4].children[i].code != structureCode["IDENTIFIER"] {
				switch s.children[4].children[i].code {
				case structureCode["L_STRING"]:
					if variable.varType != "string" {
						return createError([]string{"analyze.go", "analyze:ST_DECLARATION"}, "Excpected "+variable.varType+" got string in declaration", s.line)
					}
				case structureCode["L_BOOL"]:
					if variable.varType != "bool" {
						return createError([]string{"analyze.go", "analyze:ST_DECLARATION"}, "Excpected "+variable.varType+" got bool in declaration", s.line)
					}
				case structureCode["L_INT"]:
					valid := false
					types := []string{"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "float32", "float64"}
					for j := 0; j < len(types); j++ {
						if types[j] == variable.varType {
							valid = true
							break
						}
					}
					if !valid {
						return createError([]string{"analyze.go", "analyze:ST_DECLARATION"}, "Excpected "+variable.varType+" got int in declaration", s.line)
					}
				case structureCode["ST_CALL"]:
					var insideCall Function

					name := s.children[4].children[i].children[0].text
					var valid bool
					for i := 0; i < len(funcs); i++ {
						valid = false
						if name == funcs[i].name {
							valid = true
							insideCall = funcs[i]
							break
						}
					}
					if !valid {
						return createError([]string{"analyze.go", "analyze:ST_DECLARATION"}, "An attempt to call the non-existent function \""+name+"\" was made", s.line)
					}

					if variable.varType != insideCall.varType {
						return createError([]string{"analyze.go", "analyze:ST_DECLARATION"}, "\""+name+"\" is the wrong type, expected "+variable.varType+" got "+insideCall.varType, s.line)
					}

				default:
					return createError([]string{"analyze.go", "analyze:ST_CALL"}, "How did you even...? "+variable.name, s.line)
				}
				i += 2
				continue
			}

			i += 2
		}

	} else if s.code == structureCode["ST_FOR"] {
		n := s.children[1] // IDENTIFIER - name
		t := "int"         // IDENTIFIER - type
		v := Variable{n.text, t}
		vars = append(vars, v)
	}

	for i := 0; i < len(s.children); i++ {
		if s.children[i].code == structureCode["ST_DECLARATION"] {
			n := s.children[i].children[0] // IDENTIFIER - name
			t := s.children[i].children[2] // IDENTIFIER - type
			v := Variable{n.text, t.text}
			vars = append(vars, v)
		}

		if s.children[i].code == structureCode["ST_FUNCTION"] {
			n := s.children[i].children[1] // IDENTIFIER - name

			params := []string{}

			j := 3
			for s.children[i].children[j].code == structureCode["IDENTIFIER"] {
				t := s.children[i].children[j+2] // IDENTIFIER - type
				params = append(params, t.text)  // Add the type to the parameters of the function

				j += 4 // Next variable, otherwise this will end up on an arrow
			}

			t := s.children[i].children[len(s.children[i].children)-3]

			f := Function{n.text, params, t.text}
			funcs = append(funcs, f)
		}
		err := a.analyze(s.children[i], vars, funcs)
		if err != nil {
			return err
		}
	}

	return nil
}
