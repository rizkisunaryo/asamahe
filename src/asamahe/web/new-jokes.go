package web

import (
	//	"fmt"
	"asamahe/structs"
	"asamahe/wtpl"
	"net/http"
	"sunaryo/util/security"
	"sunaryo/util/sunhttp"
)

func HandleIndex() {
	url := "/"
	http.HandleFunc(url, index)
}

func index(w http.ResponseWriter, r *http.Request) {
	userId := "rai.roshalia"

	var req structs.NewJokesInput
	req.Viewer = userId
	pass := append([]byte(``), req.Viewer...)
	req.Key = string(security.GenKey(pass))
	req.IsBefore = 1

	url := "http://localhost:2903/api/joke/newjokes/"
	var resp structs.JokesOutput
	sunhttp.PostStruct(url, req, &resp)

	err := wtpl.WebTemplates.ExecuteTemplate(w, "index.html", resp)
}
