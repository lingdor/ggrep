package main

import (
	"bytes"
)

const (
	COLOR_NONE = "\033[0m"
	COLOR_RED  = "\033[31;5;124m"
	COLOR_BLUE = "\033[34;5;124m"
)

type ColorFormatFunc func(line string, indexes []int) string

func GrepColorRedFormat(line string, indexes []int) string {
	if len(indexes) > 1 {
		buff := &bytes.Buffer{}
		buff.WriteString(line[0:indexes[0]])
		buff.WriteString(COLOR_RED)
		buff.WriteString(line[indexes[0]:indexes[1]])
		buff.WriteString(COLOR_NONE)
		buff.WriteString(line[indexes[1]:])
		return buff.String()
	}
	return line
}
func GrepColorBLUEFormat(line string, indexes []int) string {
	if len(indexes) > 1 {
		buff := &bytes.Buffer{}
		buff.WriteString(line[0:indexes[0]])
		buff.WriteString(COLOR_BLUE)
		buff.WriteString(line[indexes[0]:indexes[1]])
		buff.WriteString(COLOR_NONE)
		buff.WriteString(line[indexes[1]:])
		return buff.String()
	}
	return line
}
