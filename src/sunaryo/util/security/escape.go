package security

import (
	//	"fmt"
	"strings"
)

func Escape(theString string) string {
	//	fmt.Printf("\n\nbefore escape:", theString)

	s := strings.Replace(theString, "<script", "&lt;script", -1)
	s = strings.Replace(s, "<style", "&lt;style", -1)
	s = strings.Replace(s, "\n", "<br />", -1)
	s = strings.Replace(s, "%", "&#37;", -1)
	s = strings.Replace(s, "\"", "\\u0022", -1)
	//	s = strings.Replace(s, "'", "&#39;", -1)

	//	fmt.Printf("\n\nafter escape:", s)

	return s
}

func EscapeForHtml(theString string) string {
	s := strings.Replace(theString, "%", "\\u0025", -1)

	return s
}
