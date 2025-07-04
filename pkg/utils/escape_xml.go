package utils

import "strings"

func EscapeXML(s string) string {
	s = strings.ReplaceAll(s, "&#xA;", "")
	return s
}
