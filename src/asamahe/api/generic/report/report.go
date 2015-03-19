package report

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

func HandleReport() {
	url := "/api/report/report/"
	http.HandleFunc(url, report)
}

func report(w http.ResponseWriter, r *http.Request) {
	var inp structs.ReportInput
	json.NewDecoder(r.Body).Decode(&inp)

	security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Reporter...))

	c := make(chan string)

	go func() {
		// Check whether the Reporter already reported the joke
		urlCheck := "http://localhost:9200/asamahe/report/_search"
		reqCheck := tpl.ReportCheckReqJsonTpl(inp.JokeId, inp.Reporter)
		var respCheck structs.CheckResponse
		err := sunhttp.Post(urlCheck, reqCheck, &respCheck)
		if err != nil {
			log.Printf("\n%s", err)
		} else if respCheck.Hits.Total > 0 {
			fmt.Fprintf(w, "{\"Status\":4}")
		} else {
			// Insert Report
			url := "http://localhost:9200/asamahe/report"
			b := reportInputJsonTemplate(inp.JokeId, inp.Reporter, inp.Long, inp.Lat)
			var respStruct structs.IdResponse
			err2 := sunhttp.Post(url, b, &respStruct)
			if err2 != nil {
				log.Printf("\n%s", err2)
			} else {
				fmt.Fprintf(w, "{\"Status\":0,\"Id\":\"%s\"}", respStruct.Id)
			}
		}

		c <- "insert Report"
	}()

	// Add and check ReportCount
	go func() {
		url2 := ("http://localhost:9200/asamahe/joke/" + inp.JokeId + "/_update")
		addReportCountReq := []byte(`{"script" : "ctx._source.ReportCount+=1"}`)
		err := sunhttp.PostNoResponse(url2, addReportCountReq)
		if err != nil {
			log.Printf("\n%s", err)
		}

		// Get the joke by JokeId
		urlGetJoke := ("http://localhost:9200/asamahe/joke/" + inp.JokeId)
		var jokeRespStruct structs.JokeResponse
		err2 := sunhttp.GetResponse(urlGetJoke, &jokeRespStruct)
		if err2 != nil {
			log.Printf("\n%s", err2)
			// If ReportCount >= 10
		} else if jokeRespStruct.Source.ReportCount >= 10 {
			// Set joke IsBlocked to 1
			urlAddScore := ("http://localhost:9200/asamahe/joke/" + inp.JokeId + "/_update")
			addScoreReq := []byte(`{"script" : "ctx._source.IsBlocked=1"}`)

			err := sunhttp.PostNoResponse(urlAddScore, addScoreReq)
			if err != nil {
				log.Printf("\n%s", err)
			}
		}
		c <- "add and check ReportCount"
	}()

	for i1 := 0; i1 < 2; i1++ {
		<-c
	}
}

func reportInputJsonTemplate(jokeId string, reporter string, long float64, lat float64) []byte {
	b := []byte(`{
  "JokeId":"`)
	b = append(b, jokeId...)
	b = append(b, []byte(`",
  "Reporter":"`)...)
	b = append(b, reporter...)
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
