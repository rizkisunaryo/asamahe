package search

import (
	"asamahe/asamahetil"
	"asamahe/structs"
	"asamahe/tpl"
	"encoding/json"
	//	"fmt"
	"net/http"
	"sunaryo/util/security"
	//	"sunaryo/util/template"
)

func HandleSearch() {
	url := "/api/search/search/"
	http.HandleFunc(url, search)
}

func search(w http.ResponseWriter, r *http.Request) {
	var inp structs.SearchInput
	json.NewDecoder(r.Body).Decode(&inp)

	if inp.Viewer != "" {
		security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Viewer...))
	}

	err := asamahetil.SaveViewActivity(inp.Viewer, "hotjokes", inp.Long, inp.Lat)
	if err != nil {
		panic(err)
	}

	// Search the jokes by Keywords
	searchJokesReq := searchJokesReqJsonTpl(inp.Keywords)
	//	fmt.Printf(string(searchJokesReq))
	tpl.JokesProcessTpl(w, inp.Viewer, searchJokesReq)
}

func searchJokesReqJsonTpl(keywords []string) []byte {
	b := []byte(`{"query":{"bool":{"should":[`)
	for i, v := range keywords {
		b = append(b, []byte(`{"wildcard":{"Content":{"value":"*`)...)
		b = append(b, v...)
		b = append(b, []byte(`*"}}}`)...)
		if i < len(keywords)-1 {
			b = append(b, []byte(`,`)...)
		}
	}
	b = append(b, []byte(`]}},"from":0,"size":100}`)...)

	//	fmt.Printf(string(b))

	return b
}
