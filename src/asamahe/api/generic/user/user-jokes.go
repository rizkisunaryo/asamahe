package user

import (
	"asamahe/asamahetil"
	"asamahe/structs"
	"asamahe/tpl"
	"encoding/json"
	//	"fmt"
	"net/http"
	"sunaryo/util/security"
	//	"sunaryo/util/template"
	"time"
)

func HandleUserJokes() {
	url := "/api/user/jokes/"
	http.HandleFunc(url, userJokes)
}

func userJokes(w http.ResponseWriter, r *http.Request) {
	var inp structs.UserJokesInput
	json.NewDecoder(r.Body).Decode(&inp)

	if inp.Viewer != "" {
		security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Viewer...))
	}

	err := asamahetil.SaveViewActivity(inp.Viewer, "jokerjokes", inp.Long, inp.Lat)
	if err != nil {
		panic(err)
	}

	// Search the jokes
	nJokesReq := searchNewJokesReqJsonTpl(inp.Joker, inp.Time, inp.IsBefore)
	//	fmt.Printf("\n%s", string(nJokesReq))
	tpl.JokerJokesProcessTpl(w, inp.Viewer, nJokesReq, inp.Joker)
}

func searchNewJokesReqJsonTpl(joker string, timestamp string, isBefore int) []byte {
	comparison, order := "gte", "asc"
	if isBefore == 1 {
		comparison, order = "lte", "desc"
	}

	b := []byte(`{"query":{"filtered":{"query":[{"match":{"IsBlocked":0}},{"match":{"Joker":"`)
	b = append(b, joker...)
	b = append(b, []byte(`"}}],"filter":{"numeric_range":{"Time":{"`)...)
	b = append(b, comparison...)
	b = append(b, []byte(`":"`)...)
	if timestamp == "" {
		tm := time.Now().Local()
		b = append(b, tm.Format("2006/01/02 15:04:05.000")...)
	} else {
		b = append(b, timestamp...)
	}
	b = append(b, []byte(`"}}}}},"sort":[{"Time":{"order":"`)...)
	b = append(b, order...)
	b = append(b, []byte(`"}}],"from":0,"size":20}`)...)

	//	fmt.Printf("%s\n", string(b))

	return b
}
