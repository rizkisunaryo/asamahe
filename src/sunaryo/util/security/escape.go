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
	s = strings.Replace(s, "\"", "\\u0022", -1)
	//	s = strings.Replace(s, "'", "&#39;", -1)

	//	fmt.Printf("\n\nafter escape:", s)

	return s
}

func Unescape(theString string) string {
	s := strings.Replace(theString, "\\u0026lt;div", "<div", -1)
	s = strings.Replace(s, "\\u0026lt;/div", "</div", -1)
	s = strings.Replace(s, "\\u0026lt;a", "<a", -1)
	s = strings.Replace(s, "\\u0026lt;/a", "</a", -1)
	s = strings.Replace(s, "\\u0026lt;img src=\\u0026#34;data", "<img src=\\u0022data", -1)
	s = strings.Replace(s, "\\u0026lt;br", "<br", -1)
	s = strings.Replace(s, "\\u0026#34;\\u003e", "\\u0022>", -1)

	//	fmt.Printf("\n\n\\u0026lt;br")
	//	fmt.Printf("\n\n%s", s)

	return s
}
