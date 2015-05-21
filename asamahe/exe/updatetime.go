package main

import (
	//	"fmt"
	//	"strconv"
	"sunaryo/util/sunhttp"
)

type Parent1 struct {
	Hits Parent11 `json:"hits"`
}

type Parent11 struct {
	Hits []Parent111 `json:"hits"`
}

type Parent111 struct {
	Id     string `json:"_id"`
	Source Joke   `json:"_source"`
}

type Joke struct {
	Joker        string
	Title        string
	Content      string
	CommentCount int
	LikeCount    int
	ReportCount  int
	Score        int
	IsAd         int
	IsBlocked    int
	Time         string
	UpdateTime   string
}

func main() {
	allJokesReq := []byte(`{
  "query": { "match_all": {} },
  "size": 1000
}`)
	var resp Parent1

	//url string, b []byte, v interface{}
	sunhttp.Post("http://localhost:9200/asamahe/joke/_search?pretty", allJokesReq, &resp)

	//	fmt.Printf(resp.Hits.Hits[0].Source.Joker)

	for _, v := range resp.Hits.Hits {
		//		fmt.Printf("\n%s", v.Source.Joker)
		//		updateTime := "2015/03/29 10:50:30.000"
		updateTime := v.Source.Time
		changeUpdateTimeUrl := "http://localhost:9200/asamahe/joke/" + v.Id + "/_update"
		changeUpdateTimeReq := []byte(`{"doc":{"UpdateTime":"`)
		changeUpdateTimeReq = append(changeUpdateTimeReq, updateTime...)
		changeUpdateTimeReq = append(changeUpdateTimeReq, []byte(`"}}`)...)
		/*err3 :=*/ sunhttp.PostNoResponse(changeUpdateTimeUrl, changeUpdateTimeReq)
		//		if err3 != nil {
		//			log.Printf("\n%s", err3)
		//		}
		//		sunhttp.PostNoResponse("http://localhost:9200/asamahe/joke/"+v.Id, updateReq)
	}
}
