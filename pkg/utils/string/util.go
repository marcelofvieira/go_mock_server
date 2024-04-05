package string

import (
	"strings"
)

func ReplaceTabsToSpaces(str string) string {
	return strings.ReplaceAll(str, "\t", " ")
}

func ReplaceNewLinesToSpaces(str string) string {
	str = strings.ReplaceAll(str, "\\n", "")
	return strings.ReplaceAll(str, "\n", " ")
}

func RemoveSpaces(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

func RemoveParenthesis(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "(", "\\("), ")", "\\)")
}

func BetweenPosition(value string, a int, b int) string {
	return value[a:b]
}

func Between(value string, a string, b string) (string, int, int) {
	posFirst := strings.Index(value, a)

	if posFirst == -1 {
		return "", -1, -1
	}

	posLast := strings.Index(value, b)

	if posLast == -1 {
		return "", -1, -1
	}

	posFirstAdjusted := posFirst + len(a)

	if posFirstAdjusted >= posLast {
		return "", -1, -1
	}

	return value[posFirstAdjusted:posLast], posFirstAdjusted, posLast
}

func before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}
