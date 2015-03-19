package security

import (
	"strings"
)

func Escape(theString string) string {
	s := strings.Replace(theString, "<", "&lt;", -1)
	s = strings.Replace(s, "\n", "<br />", -1)
	s = strings.Replace(s, "\"", "&#34;", -1)
	s = strings.Replace(s, "'", "&#39;", -1)
	return s
}
