package user

import (
	"asamahe/asamahetil"
	"asamahe/structs"
	"asamahe/tpl"
	"encoding/json"
	"fmt"
	"net/http"
	"sunaryo/util/security"
	"sunaryo/util/sunhttp"
	//	"sunaryo/util/template"
	//	"log"
)

func HandleTopJokers() {
	url := "/api/user/topjokers/"
	http.HandleFunc(url, topJokers)
}

func topJokers(w http.ResponseWriter, r *http.Request) {
	var inp structs.TopJokersInput
	json.NewDecoder(r.Body).Decode(&inp)

	if inp.Viewer != "" {
		security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Viewer...))
	}

	err := asamahetil.SaveViewActivity(inp.Viewer, "topjokers", inp.Long, inp.Lat)
	if err != nil {
		panic(err)
	}

	// Search the jokes
	topJokersReq := []byte(`{"aggs":{"group_by_joker":{"terms":{"field":"Joker","order":{"sumScore":"desc"}},"aggs":{"sumScore":{"sum":{"field":"Score"}}}}},"from":0,"size":100}`)
	topJokersUrl := "http://localhost:9200/asamahe/joke/_search"
	var jqr structs.TopJokersResponse
	err2 := sunhttp.Post(topJokersUrl, topJokersReq, &jqr)
	if err2 != nil {
		panic(err2)
	} else {

		var jokersOutput structs.JokersOutput
		for _, v := range jqr.Aggregations.Group.JokerScores {
			var jokerResp structs.UserResponse
			tpl.JokerProcessTpl(w, &jokerResp, v.Key)
			jokersOutput.Jokers = append(jokersOutput.Jokers, jokerResp.Source)
			jokersOutput.Jokers[len(jokersOutput.Jokers)-1].Id = v.Key
		}

		jsStr, _ := json.Marshal(jokersOutput)
		fmt.Fprintf(w, string(jsStr))
	}
}
