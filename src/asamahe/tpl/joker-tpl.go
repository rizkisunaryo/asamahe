package tpl

import (
	//	"asamahe/structs"
	"net/http"
	"sunaryo/util/sunhttp"
	//	"sunaryo/util/template"
)

func JokerProcessTpl(w http.ResponseWriter, v interface{}, id string) {
	jokerUrl := ("http://localhost:9200/asamahe/user/" + id)
	err := sunhttp.GetResponse(jokerUrl, v)
	if err != nil {
		panic(err)
	}
}
