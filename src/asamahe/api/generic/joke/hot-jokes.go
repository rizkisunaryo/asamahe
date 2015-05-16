package joke

import (
	"asamahe/asamahetil"
	"asamahe/structs"
	"asamahe/tpl"
	"encoding/json"
	"net/http"
	"strconv"
	"sunaryo/util/security"
	//	"sunaryo/util/template"
)

func HandleHotJokes() {
	url := "/api/joke/hotjokes/"
	http.HandleFunc(url, hotJokes)
}

func hotJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var inp structs.HotJokesInput
	json.NewDecoder(r.Body).Decode(&inp)

	if inp.Viewer != "" {
		security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Viewer...))
	}

	err := asamahetil.SaveViewActivity(inp.Viewer, "hotjokes", inp.Long, inp.Lat)
	if err != nil {
		panic(err)
	}

	// Search the jokes
	nJokesReq := searchHotJokesReqJsonTpl(inp.Score, inp.IsBefore, inp.From, inp.Size)
	tpl.JokesProcessTpl(w, inp.Viewer, nJokesReq)
}

func searchHotJokesReqJsonTpl(score int, isBefore int, from int, size int) []byte {
	comparison, order := "gte", "asc"
	if isBefore == 1 {
		comparison, order = "lte", "desc"
	}

	b := []byte(`{"`)
	if score != -1000 {
		b = append(b, []byte(`query":{"filtered":{"query":{"match":{"IsBlocked":0}},"filter":{"numeric_range":{"Score":{"`)...)
		b = append(b, comparison...)
		b = append(b, []byte(`":`)...)
		b = append(b, strconv.Itoa(score)...)
		b = append(b, []byte(`}}}}},"`)...)
	}
	b = append(b, []byte(`sort":[{"Score":{"order":"`)...)
	b = append(b, order...)
	b = append(b, []byte(`"}}],"from":`)...)
	b = append(b, strconv.Itoa(from)...)
	b = append(b, []byte(`,"size":`)...)
	b = append(b, strconv.Itoa(size)...)
	b = append(b, []byte(`}`)...)

	//	fmt.Printf("%s\n", string(b))

	return b
}
