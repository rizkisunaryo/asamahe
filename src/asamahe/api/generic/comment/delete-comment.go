package comment

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

func HandleDeleteComment() {
	url := "/api/comment/deletecomment/"
	http.HandleFunc(url, deleteComment)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var inp structs.DeleteCommentInput
	json.NewDecoder(r.Body).Decode(&inp)

	security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.UserId...))

	c := make(chan string)

	// Delete comment
	go func() {
		// Get the Comment
		urlGetComment := ("http://localhost:9200/asamahe/comment/" + inp.CommentId)
		var commentRespStruct structs.CommentResponse
		err2 := sunhttp.GetResponse(urlGetComment, &commentRespStruct)
		if err2 != nil {
			log.Printf("\n%s", err2)
		} else {
			// Check whether UserId is the Commentator
			if inp.UserId == commentRespStruct.Source.Commentator {
				// If UserId is the Commentator,
				// he/she can delete his/her comment
				url := ("http://localhost:9200/asamahe/comment/" + inp.CommentId)
				err3 := sunhttp.Delete(url)
				if err3 != nil {
					log.Printf("\n%s", err3)
				} else {
					fmt.Fprintf(w, "{\"Status\":0}")
				}
			}
		}

		c <- "delete comment"
	}()

	// Reduce CommentCount and Score
	go func() {
		url2 := ("http://localhost:9200/asamahe/joke/" + inp.JokeId + "/_update")
		addCommentCountReq := []byte(`{"script" : "ctx._source.CommentCount-=1"}`)
		err := sunhttp.PostNoResponse(url2, addCommentCountReq)
		if err != nil {
			log.Printf("\n%s", err)
		}

		// Get the joke by JokeId
		urlGetJoke := ("http://localhost:9200/asamahe/joke/" + inp.JokeId)
		var jokeRespStruct structs.JokeResponse
		err2 := sunhttp.GetResponse(urlGetJoke, &jokeRespStruct)
		if err2 != nil {
			log.Printf("\n%s", err2)
			// Check whether UserId is not the Joker
		} else if jokeRespStruct.Source.Joker != inp.UserId {
			// If not, gonna reduce the score by 3
			urlAddScore := ("http://localhost:9200/asamahe/joke/" + inp.JokeId + "/_update")
			addScoreReq := []byte(`{"script" : "ctx._source.Score-=3"}`)

			err := sunhttp.PostNoResponse(urlAddScore, addScoreReq)
			if err != nil {
				log.Printf("\n%s", err)
			}
		}

		c <- "reduce CommentCount and Score"
	}()

	for i1 := 0; i1 < 2; i1++ {
		<-c
	}
}
