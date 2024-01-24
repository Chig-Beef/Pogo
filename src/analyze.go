package main

import "errors"

type Analyzer struct {
}

func (a *Analyzer) analyze(s Structure, vars []Variable) error {
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
			return errors.New("[Analyzer (analyze-manipulation)] An attempt to manipulate an uninitialized variable was made")
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
				return errors.New("[Analyzer (analyze-expression)] An uninitialized variable was used in an expression")
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
				return errors.New("[Analyzer (analyze-coparison)] An uninitialized variable was used in a comparison")
			}
		}
	} else if s.code == structureCode["ST_CALL"] {
		if s.children[2].code == structureCode["IDENTIFIER"] {
			name := s.children[2].text
			var valid bool
			for i := 0; i < len(vars); i++ {
				valid = false
				if name == vars[i].name {
					valid = true
					break
				}
			}
			if !valid {
				return errors.New("[Analyzer (analyze-call)] An uninitialized variable was used in a function call")
			}
		}
	}

	for i := 0; i < len(s.children); i++ {
		err := a.analyze(s.children[i], vars)
		if err != nil {
			return err
		}

		if s.children[i].code == structureCode["ST_DECLARATION"] {
			n := s.children[i].children[0] // IDENTIFIER - name
			t := s.children[i].children[2] // IDENTIFIER - type
			v := Variable{n.text, t.text}
			vars = append(vars, v)
		}
	}

	return nil
}

type Variable struct {
	name    string
	varType string
}
