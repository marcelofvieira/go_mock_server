package stringutils

import "strings"

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
