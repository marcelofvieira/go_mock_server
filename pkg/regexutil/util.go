package regexutil

import "regexp"

const (
	FindResponseVariablePattern = "[\\$]?[\\$(NB)]{(.+?)}"
	FindVariableContextPattern  = "[\\$]?[\\$(NB)]{([a-zA-Z0-9]+)\\."
	//FindVariablePattern         = "\\{([^{}]+)\\}"
	//FindBodyVariablePattern     = "\\{([^{}]+)\\}"
	FindVariablePattern      = "\\{{(.+?)\\}}"
	FindBodyVariablePattern  = "\\{{(.+?)\\}}"
	FindVariableValuePattern = "([^/]+)"
	FindToFinalPattern       = "+$"
	//Locate to variables to work with number or boolean
	FindNumberBooleanVariablePattern = "[(NB)]{"
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
