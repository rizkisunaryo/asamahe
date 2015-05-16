package joke

import (
	"encoding/json"
	//	"fmt"
	"asamahe/asamahetil"
	"asamahe/structs"
	"asamahe/tpl"
	//	"io/ioutil"
	"net/http"
	"sunaryo/util/security"
	//	"sunaryo/util/template"
	"time"
)

func HandleNewJokes() {
	url := "/api/joke/newjokes/"
	http.HandleFunc(url, newJokes)
}

func newJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//	body, _ := ioutil.ReadAll(r.Body)
	//	fmt.Printf(string(body))

	var inp structs.NewJokesInput
	json.NewDecoder(r.Body).Decode(&inp)

	if inp.Viewer != "" {
		security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Viewer...))
	}

	err := asamahetil.SaveViewActivity(inp.Viewer, "newjokes", inp.Long, inp.Lat)
	if err != nil {
		panic(err)
	}

	// Search the jokes
	nJokesReq := searchNewJokesReqJsonTpl(inp.Time, inp.IsBefore)
	tpl.JokesProcessTpl(w, inp.Viewer, nJokesReq)
}

func searchNewJokesReqJsonTpl(timestamp string, isBefore int) []byte {
	comparison, order := "gte", "asc"
	if isBefore == 1 {
		comparison, order = "lte", "desc"
	}

	b := []byte(`{"query":{"filtered":{"query":{"match":{"IsBlocked":0}},"filter":{"numeric_range":{"UpdateTime":{"`)
	b = append(b, comparison...)
	b = append(b, []byte(`":"`)...)
	if timestamp == "" {
		tm := time.Now().Local()
		b = append(b, tm.Format("2006/01/02 15:04:05.000")...)
	} else {
		b = append(b, timestamp...)
	}
	b = append(b, []byte(`"}}}}},"sort":[{"UpdateTime":{"order":"`)...)
	b = append(b, order...)
	b = append(b, []byte(`"}}],"from":0,"size":20}`)...)

	//	fmt.Printf("%s\n", string(b))

	return b
}
