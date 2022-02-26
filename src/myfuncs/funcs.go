package myfuncs

import "strings"

func PrimeraMayus(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}

	s = strings.ToUpper(string(s[0])) + s[1:]
	return s
}
