package like

import (
	"asamahe/structs"
	"encoding/json"
	"fmt"
	"net/http"
	"sunaryo/util/security"
	"sunaryo/util/sunhttp"
	//	"sunaryo/util/template"
	"log"
)

func HandleDislike() {
	url := "/api/like/dislike/"
	http.HandleFunc(url, dislike)
}

func dislike(w http.ResponseWriter, r *http.Request) {
	var inp structs.LikeInput
	json.NewDecoder(r.Body).Decode(&inp)

	security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Liker...))

	c := make(chan string)

	// Dislike
	go func() {
		// Delete like
		delLikeUrl := "http://localhost:9200/asamahe/like/_query"
		delLikeReq := delLikeReqJsonTpl(inp.Liker, inp.JokeId)
		err := sunhttp.DeleteStruct(delLikeUrl, delLikeReq)
		if err != nil {
			log.Printf("\n%s", err)
		} else {
			fmt.Fprintf(w, "{\"Status\":0}")
		}
		c <- "dislike"
	}()

	// Reduce LikeCount and Score
	go func() {
		url2 := ("http://localhost:9200/asamahe/joke/" + inp.JokeId + "/_update")
		addLikeCountReq := []byte(`{"script" : "ctx._source.LikeCount-=1"}`)
		err := sunhttp.PostNoResponse(url2, addLikeCountReq)
		if err != nil {
			log.Printf("\n%s", err)
		}

		// Get the joke, to get the Joker
		urlGetJoke := ("http://localhost:9200/asamahe/joke/" + inp.JokeId)
		var jokeRespStruct structs.JokeResponse
		err2 := sunhttp.GetResponse(urlGetJoke, &jokeRespStruct)
		if err2 != nil {
			log.Printf("\n%s", err2)
			// If Joker is not Liker himself
		} else if jokeRespStruct.Source.Joker != inp.Liker {
			// Reduce Score by 1
			urlAddScore := ("http://localhost:9200/asamahe/joke/" + inp.JokeId + "/_update")
			addScoreReq := []byte(`{"script" : "ctx._source.Score-=1"}`)

			err := sunhttp.PostNoResponse(urlAddScore, addScoreReq)
			if err != nil {
				log.Printf("\n%s", err)
			}
		}
		c <- "reduce LikeCount and Score"
	}()

	for i1 := 0; i1 < 2; i1++ {
		<-c
	}
}

func delLikeReqJsonTpl(liker, jokeId string) []byte {
	b := []byte(`{"query":{"bool":{"must":[{"match":{"Liker":"`)
	b = append(b, liker...)
	b = append(b, []byte(`"}},{"match":{"JokeId":"`)...)
	b = append(b, jokeId...)
	b = append(b, []byte(`"}}]}}}`)...)
	return b
}
