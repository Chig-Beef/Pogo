package main

import (
	"errors"
	"strconv"
)

func createError(funcLine []string, message string, line int) error {
	output := ""
	for i := 0; i < len(funcLine); i++ {
		output += funcLine[i]
		output += " -> "
	}
	output += "line: " + strconv.Itoa(line) + "\n"
	output += message
	return errors.New(output)
}
