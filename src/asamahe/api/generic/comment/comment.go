package comment

import (
	"asamahe/structs"
	"encoding/json"
	"fmt"
	"net/http"
	"sunaryo/util"
	"sunaryo/util/security"
	"sunaryo/util/sunhttp"
	//	"sunaryo/util/template"
	"log"
	"time"
)

func HandleComment() {
	url := "/api/comment/comment/"
	http.HandleFunc(url, comment)
}

func comment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var t structs.CommentInput
	json.NewDecoder(r.Body).Decode(&t)

	if t.Commentator != "" && t.Key != "" && t.JokeId != "" && t.Comment != "" {
		security.CheckKey(append([]byte(``), t.Key...), append([]byte(``), t.Commentator...))

		c := make(chan string)

		// Insert comment
		go func() {
			url := "http://localhost:9200/asamahe/comment"
			b := commentInputJsonTemplate(t.JokeId, security.Escape(t.Commentator), security.Escape(t.Comment), t.Long, t.Lat)
			var respStruct structs.IdResponse
			err := sunhttp.Post(url, b, &respStruct)
			if err != nil {
				log.Printf("\n%s", err)
			} else {
				fmt.Fprintf(w, "{\"Status\":0,\"Id\":\"%s\"}", respStruct.Id)
			}
			c <- "insert comment"
		}()

		// Add CommentCount and Score
		go func() {
			url2 := ("http://localhost:9200/asamahe/joke/" + t.JokeId + "/_update")
			addCommentCountReq := []byte(`{"script" : "ctx._source.CommentCount+=1"}`)
			err := sunhttp.PostNoResponse(url2, addCommentCountReq)
			if err != nil {
				log.Printf("\n%s", err)
			}

			// Get Joke
			// The intention is to get the Joker
			urlGetJoke := ("http://localhost:9200/asamahe/joke/" + t.JokeId)
			var jokeRespStruct structs.JokeResponse
			err2 := sunhttp.GetResponse(urlGetJoke, &jokeRespStruct)
			if err2 != nil {
				log.Printf("\n%s", err2)
				// Check whether the Commentator is not the Joker itself
			} else if jokeRespStruct.Source.Joker != t.Commentator {
				// If not the Joker itself,
				// add Score by 3
				urlAddScore := ("http://localhost:9200/asamahe/joke/" + t.JokeId + "/_update")
				addScoreReq := []byte(`{"script" : "ctx._source.Score+=3"}`)

				err := sunhttp.PostNoResponse(urlAddScore, addScoreReq)
				if err != nil {
					log.Printf("\n%s", err)
				}
			}
			c <- "Add CommentCount and Score"
		}()

		for i1 := 0; i1 < 2; i1++ {
			//		fmt.Printf("%s\n", <-c)
			<-c
		}
	} else {
		fmt.Fprintf(w, "{\"Status\":1}")
	}
}

func commentInputJsonTemplate(jokeId string, commentator string, comment string, long float64, lat float64) []byte {
	b := []byte(`{
  "JokeId":"`)
	b = append(b, jokeId...)
	b = append(b, []byte(`",
  "Commentator":"`)...)
	b = append(b, commentator...)
	b = append(b, []byte(`",
  "Comment":"`)...)
	b = append(b, comment...)
	b = append(b, []byte(`",
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
