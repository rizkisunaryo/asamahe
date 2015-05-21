package tpl

import (
	"asamahe/structs"
	"encoding/json"
	"fmt"
	"net/http"
	"sunaryo/util/sunhttp"
	//	"sunaryo/util/template"
	"log"
	"sunaryo/util/security"
)

func JokesProcessTpl(w http.ResponseWriter, viewer string, b []byte) {
	//	var nJokesOutput structs.JokesOutput
	nJokesOutput := JokesGetFromESProcessTpl(w, viewer, b)
	jsStr, _ := json.Marshal(nJokesOutput)
	fmt.Fprintf(w, security.EscapeForHtml(string(jsStr)))
}

func JokerJokesProcessTpl(w http.ResponseWriter, viewer string, b []byte, joker string) {
	var jokerJokes structs.JokerJokesOutput

	c := make(chan string)

	// Get Joker
	go func() {
		urlGetJoker := ("http://localhost:9200/asamahe/user/" + joker)
		var jokerResp structs.UserResponse
		err := sunhttp.GetResponse(urlGetJoker, &jokerResp)
		if err != nil {
			log.Printf("\n%s", err)
		} else {
			jokerResp.Source.Id = joker
			jokerJokes.Joker = jokerResp.Source
		}
		c <- "get joker"
	}()

	// Get jokes by Joker
	go func() {
		//		var nJokesOutput structs.JokesOutput
		nJokesOutput := JokesGetFromESProcessTpl(w, viewer, b)
		jokerJokes.Jokes = nJokesOutput.Jokes
		c <- "get jokes"
	}()

	for i := 0; i < 2; i++ {
		<-c
	}

	jsStr, _ := json.Marshal(jokerJokes)
	fmt.Fprintf(w, security.EscapeForHtml(string(jsStr)))
}

func JokesGetFromESProcessTpl(w http.ResponseWriter, viewer string, b []byte) structs.JokesOutput {
	var nJokesOutput structs.JokesOutput

	nJokesUrl := "http://localhost:9200/asamahe/joke/_search"
	var jqr structs.JokeQueryResponse
	err2 := sunhttp.Post(nJokesUrl, b, &jqr)
	if err2 != nil {
		panic(err2)
	} else {

		for _, v := range jqr.Hits.Hits {
			// Add searched joke to nJokesOutput.Jokes
			nJokesOutput.Jokes = append(nJokesOutput.Jokes, v.Source)
			// set JokeId of nJokesOutput.Jokes element
			nJokesOutput.Jokes[len(nJokesOutput.Jokes)-1].JokeId = v.Id

			// Set JokerName and JokerPicUrl
			var jokerResp structs.UserResponse
			JokerProcessTpl(w, &jokerResp, nJokesOutput.Jokes[len(nJokesOutput.Jokes)-1].Joker)
			nJokesOutput.Jokes[len(nJokesOutput.Jokes)-1].JokerName = jokerResp.Source.Name
			nJokesOutput.Jokes[len(nJokesOutput.Jokes)-1].JokerPicUrl = jokerResp.Source.PicUrl

			// Check whether the joke IsLiked
			if viewer == "" {
				nJokesOutput.Jokes[len(nJokesOutput.Jokes)-1].IsLiked = 0
			} else {
				likeCheckUrl := "http://localhost:9200/asamahe/like/_search"
				likeCheckReq := LikeCheckReqJsonTpl(v.Id, viewer)
				var likeCheckResp structs.CheckResponse
				err := sunhttp.Post(likeCheckUrl, likeCheckReq, &likeCheckResp)
				if err != nil {
					log.Printf("\n%s", err)
				} else if likeCheckResp.Hits.Total > 0 {
					nJokesOutput.Jokes[len(nJokesOutput.Jokes)-1].IsLiked = 1
				}
			}

			// Check whether the joke IsReported
			if viewer == "" {
				nJokesOutput.Jokes[len(nJokesOutput.Jokes)-1].IsReported = 0
			} else {
				reportCheckUrl := "http://localhost:9200/asamahe/report/_search"
				reportCheckReq := ReportCheckReqJsonTpl(v.Id, viewer)
				var reportCheckResp structs.CheckResponse
				err := sunhttp.Post(reportCheckUrl, reportCheckReq, &reportCheckResp)
				if err != nil {
					log.Printf("\n%s", err)
				} else if reportCheckResp.Hits.Total > 0 {
					nJokesOutput.Jokes[len(nJokesOutput.Jokes)-1].IsReported = 1
				}
			}
		}
	}

	return nJokesOutput
}
