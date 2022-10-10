package utils

import "strings"

func ConcatStr(s ...string) string {
	var sb strings.Builder

	for _, a := range s {
		sb.WriteString(a)
	}

	return sb.String()
}
