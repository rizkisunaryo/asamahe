package like

import (
	"asamahe/structs"
	"asamahe/tpl"
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

func HandleLike() {
	url := "/api/like/like/"
	http.HandleFunc(url, like)
}

func like(w http.ResponseWriter, r *http.Request) {
	var inp structs.LikeInput
	json.NewDecoder(r.Body).Decode(&inp)

	security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Liker...))

	c := make(chan string)

	go func() {
		// Check whether the Liker already liked the joke
		urlCheck := "http://localhost:9200/asamahe/like/_search"
		reqCheck := tpl.LikeCheckReqJsonTpl(inp.JokeId, inp.Liker)
		var respCheck structs.CheckResponse
		err := sunhttp.Post(urlCheck, reqCheck, &respCheck)
		if err != nil {
			log.Printf("\n%s", err)
		} else if respCheck.Hits.Total > 0 {
			fmt.Fprintf(w, "{\"Status\":4}")
		} else {
			// Insert like
			url := "http://localhost:9200/asamahe/like"
			b := likeInputJsonTemplate(inp.JokeId, inp.Liker, inp.Long, inp.Lat)
			var respStruct structs.IdResponse
			err := sunhttp.Post(url, b, &respStruct)
			if err != nil {
				log.Printf("\n%s", err)
			} else {
				fmt.Fprintf(w, "{\"Status\":0,\"Id\":\"%s\"}", respStruct.Id)
			}
		}
		c <- "insert like"
	}()

	// Add LikeCount and Score
	go func() {
		url2 := ("http://localhost:9200/asamahe/joke/" + inp.JokeId + "/_update")
		addLikeCountReq := []byte(`{"script" : "ctx._source.LikeCount+=1"}`)
		err := sunhttp.PostNoResponse(url2, addLikeCountReq)
		if err != nil {
			log.Printf("\n%s", err)
		}

		// Get the joke by JokeId
		urlGetJoke := ("http://localhost:9200/asamahe/joke/" + inp.JokeId)
		var jokeRespStruct structs.JokeResponse
		err2 := sunhttp.GetResponse(urlGetJoke, &jokeRespStruct)
		if err2 != nil {
			log.Printf("\n%s", err2)
			// If Liker is not the Joker
		} else if jokeRespStruct.Source.Joker != inp.Liker {
			// Add Score by 1
			urlAddScore := ("http://localhost:9200/asamahe/joke/" + inp.JokeId + "/_update")
			addScoreReq := []byte(`{"script" : "ctx._source.Score+=1"}`)

			err := sunhttp.PostNoResponse(urlAddScore, addScoreReq)
			if err != nil {
				log.Printf("\n%s", err)
			}
		}
		c <- "add LikeCount and Score"
	}()

	for i1 := 0; i1 < 2; i1++ {
		<-c
	}
}

func likeInputJsonTemplate(jokeId string, liker string, long float64, lat float64) []byte {
	b := []byte(`{
  "JokeId":"`)
	b = append(b, jokeId...)
	b = append(b, []byte(`",
  "Liker":"`)...)
	b = append(b, liker...)
	b = append(b, []byte(`","Time":"`)...)
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
