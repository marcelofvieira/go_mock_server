package regexutil

import "regexp"

const (
	FindVariablePattern      = "\\{{(.+?)\\}}"
	FindBodyVariablePattern  = "\\{{(.+?)\\}}"
	FindVariableValuePattern = "([^/]+)"
	FindToFinalPattern       = "+$"
)

func FindStringRegex(pattern, text string) bool {
	validRegex := regexp.MustCompile(pattern)

	return validRegex.MatchString(text)
}

func FindStringValuesRegex(pattern, text string) (bool, [][]string) {
	validRegex := regexp.MustCompile(pattern)

	matches := validRegex.FindAllStringSubmatch(text, -1)

	return len(matches) > 0, matches
}
