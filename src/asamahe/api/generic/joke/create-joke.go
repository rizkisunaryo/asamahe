package joke

import (
	"asamahe/structs"
	"encoding/json"
	"fmt"
	"net/http"
	"sunaryo/util"
	"sunaryo/util/security"
	"sunaryo/util/sunhttp"
	"sunaryo/util/template"
	"time"
)

func HandleCreateJoke() {
	url := "/api/joke/createjoke/"
	http.HandleFunc(url, createJoke)
}

func createJoke(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var t structs.CreateJokeInput
	json.NewDecoder(r.Body).Decode(&t)

	if t.Joker != "" && t.Key != "" && t.Content != "" && len(t.Content) <= 163840 && len(t.Title) <= 1024 {
		security.CheckKey(append([]byte(``), t.Key...), append([]byte(``), t.Joker...))

		url := "http://localhost:9200/asamahe/joke"
		// fmt.Printf(security.Escape(t.Content))
		b := createJokeInputJsonTemplate(t.Joker, security.Escape(t.Title), security.Escape(t.Content), t.Long, t.Lat)
		var respStruct structs.CreateJokeResponse
		succFn := func() {
			fmt.Fprintf(w, "{\"Status\":0,\"Id\":\"%s\"}", respStruct.Id)
		}
		err := sunhttp.Post(url, b, &respStruct)
		template.IfError(err, w, 1, succFn)
	} else {
		fmt.Fprintf(w, "{\"Status\":1}")
	}
}

func createJokeInputJsonTemplate(joker string, title string, content string, long float64, lat float64) []byte {
	b := []byte(`{
  "Joker": "`)
	b = append(b, joker...)
	b = append(b, []byte(`", 
        "Title": "`)...)
	b = append(b, title...)
	b = append(b, []byte(`", 
        "Content":"`)...)
	b = append(b, content...)
	b = append(b, []byte(`",
        "CommentCount":0,
        "LikeCount":0,
        "ReportCount":0,
        "Score":0,
        "IsAd":0,
        "IsBlocked":0,
        "Time":"`)...)
	tm := time.Now().Local()
	b = append(b, tm.Format("2006/01/02 15:04:05.000")...)
	if long != 0.0 || lat != 0.0 {
		b = append(b, []byte(`","Geo":{"Loc":[`)...)
		b = append(b, util.FloatToString(long)...)
		b = append(b, ","...)
		b = append(b, util.FloatToString(lat)...)
		b = append(b, "]}}"...)
	} else {
		b = append(b, []byte(`"}`)...)
	}

	return b
}
