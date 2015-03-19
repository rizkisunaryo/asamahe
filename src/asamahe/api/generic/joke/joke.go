package joke

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
	//	"io/ioutil"
	"log"
)

func HandleJoke() {
	url := "/api/joke/joke/"
	http.HandleFunc(url, joke)
}

func joke(w http.ResponseWriter, r *http.Request) {
	//	body, _ := ioutil.ReadAll(r.Body)
	//	fmt.Printf(string(body))

	var inp structs.JokeInput
	json.NewDecoder(r.Body).Decode(&inp)

	if inp.Viewer != "" {
		security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Viewer...))
	}

	var jokeOutput structs.JokeOutput
	jokeOutput.Status = 0

	//Get the Joke by JokeId
	urlGetJoke := ("http://localhost:9200/asamahe/joke/" + inp.JokeId)
	var jokeResp structs.JokeResponse
	err := sunhttp.GetResponse(urlGetJoke, &jokeResp)
	if err != nil {
		panic(err)
	}

	if jokeResp.Source.IsBlocked == 1 {
		fmt.Fprintf(w, "{\"Status\":4}")
		return
	}
	jokeResp.Source.JokeId = inp.JokeId

	c := make(chan string)

	// Save the activity as view activity
	go func() {
		err := asamahetil.SaveViewActivity(inp.Viewer, "joke", inp.Long, inp.Lat)
		if err != nil {
			log.Printf("\n%s", err)
		}
		c <- "viewact"
	}()

	//	To get the JokerName from user category/table
	go func() {
		jokerUrl := ("http://localhost:9200/asamahe/user/" + jokeResp.Source.Joker)
		var jokerResp structs.UserResponse
		err := sunhttp.GetResponse(jokerUrl, &jokerResp)
		if err != nil {
			log.Printf("\n%s", err)
		} else {
			jokeResp.Source.JokerName = jokerResp.Source.Name
			jokeResp.Source.JokerPicUrl = jokerResp.Source.PicUrl
			jokeOutput.Joke = jokeResp.Source
		}
		c <- "joker"
	}()

	// Whether IsLiked by Viewer
	go func() {
		likeCheckUrl := "http://localhost:9200/asamahe/like/_search"
		likeCheckReq := tpl.LikeCheckReqJsonTpl(inp.JokeId, inp.Viewer)
		var likeCheckResp structs.CheckResponse
		err := sunhttp.Post(likeCheckUrl, likeCheckReq, &likeCheckResp)
		if err != nil {
			log.Printf("\n%s", err)
		} else {
			if likeCheckResp.Hits.Total > 0 {
				jokeResp.Source.IsLiked = 1
			}
			jokeOutput.Joke = jokeResp.Source
		}
		c <- "IsLiked"
	}()

	// Whether IsReported by Viewer
	go func() {
		reportCheckUrl := "http://localhost:9200/asamahe/report/_search"
		reportCheckReq := tpl.ReportCheckReqJsonTpl(inp.JokeId, inp.Viewer)
		var reportCheckResp structs.CheckResponse
		err := sunhttp.Post(reportCheckUrl, reportCheckReq, &reportCheckResp)
		if err != nil {
			log.Printf("\n%s", err)
		} else {
			if reportCheckResp.Hits.Total > 0 {
				jokeResp.Source.IsReported = 1
			}
			jokeOutput.Joke = jokeResp.Source
		}
		c <- "IsReported"
	}()

	// Find all comments of the joke
	go func() {
		commentUrl := "http://localhost:9200/asamahe/comment/_search"
		commentReq := searcCommentshByJokeIdReqJson(inp.JokeId)
		var cqr structs.CommentQueryResponse
		err := sunhttp.Post(commentUrl, commentReq, &cqr)
		if err != nil {
			log.Printf("\n%s", err)
		} else {

			for _, v := range cqr.Hits.Hits {
				if v.Source.IsBlocked == 0 {
					jokeOutput.Comments = append(jokeOutput.Comments, v.Source)

					commentatorUrl := ("http://localhost:9200/asamahe/user/" + v.Source.Commentator)
					var commentatorResp structs.UserResponse
					err := sunhttp.GetResponse(commentatorUrl, &commentatorResp)
					if err != nil {
						log.Printf("\n%s", err)
					} else {
						jokeOutput.Comments[len(jokeOutput.Comments)-1].CommentId = v.Id
						jokeOutput.Comments[len(jokeOutput.Comments)-1].CommentatorName = commentatorResp.Source.Name
						jokeOutput.Comments[len(jokeOutput.Comments)-1].CommentatorPicUrl = commentatorResp.Source.PicUrl
					}
				}
			}
		}
		c <- "comments"
	}()

	// Find all likes of the joke
	go func() {
		likeUrl := "http://localhost:9200/asamahe/like/_search"
		likeReq := searchByJokeIdReqJson(inp.JokeId)
		var lqr structs.LikeQueryResponse
		err := sunhttp.Post(likeUrl, likeReq, &lqr)
		if err != nil {
			log.Printf("\n%s", err)
		} else {
			for _, v := range lqr.Hits.Hits {
				jokeOutput.Likes = append(jokeOutput.Likes, v.Source)

				likerUrl := ("http://localhost:9200/asamahe/user/" + v.Source.Liker)
				var likerResp structs.UserResponse
				err := sunhttp.GetResponse(likerUrl, &likerResp)
				if err != nil {
					log.Printf("\n%s", err)
				} else {
					jokeOutput.Likes[len(jokeOutput.Likes)-1].LikerName = likerResp.Source.Name
				}
			}
		}
		c <- "likes"
	}()

	// Find all reports of the joke
	go func() {
		reportUrl := "http://localhost:9200/asamahe/report/_search"
		reportReq := searchByJokeIdReqJson(inp.JokeId)
		var rqr structs.ReportQueryResponse
		err := sunhttp.Post(reportUrl, reportReq, &rqr)
		if err != nil {
			log.Printf("\n%s", err)
		} else {
			for _, v := range rqr.Hits.Hits {
				jokeOutput.Reports = append(jokeOutput.Reports, v.Source)

				reportrUrl := ("http://localhost:9200/asamahe/user/" + v.Source.Reporter)
				var reportrResp structs.UserResponse
				err := sunhttp.GetResponse(reportrUrl, &reportrResp)
				if err != nil {
					log.Printf("\n%s", err)
				} else {
					jokeOutput.Reports[len(jokeOutput.Reports)-1].ReporterName = reportrResp.Source.Name
				}
			}
		}
		c <- "reports"
	}()

	for i1 := 0; i1 < 7; i1++ {
		//		fmt.Printf("%s\n", <-c)
		<-c
	}

	jsStr, _ := json.Marshal(jokeOutput)
	fmt.Fprintf(w, string(jsStr))
}

func searchByJokeIdReqJson(jokeId string) []byte {
	b := []byte(`{"query":{"bool":{"must":[{"match":{"JokeId":"`)
	b = append(b, jokeId...)
	b = append(b, []byte(`"}}]}}}`)...)
	return b
}

func searcCommentshByJokeIdReqJson(jokeId string) []byte {
	b := []byte(`{"query":{"bool":{"must":[{"match":{"JokeId":"`)
	b = append(b, jokeId...)
	b = append(b, []byte(`"}}]}},"sort":[{"Time":{"order":"asc"}}],"from":0,"size": 100}`)...)
	return b
}
