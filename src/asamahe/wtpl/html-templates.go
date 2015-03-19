package wtpl

import (
	"html/template"
)

var WebTemplates = template.Must(template.ParseFiles("index.html", "newjokes.html"))
