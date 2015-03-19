package report

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

func HandleUnreport() {
	url := "/api/report/unreport/"
	http.HandleFunc(url, unreport)
}

func unreport(w http.ResponseWriter, r *http.Request) {
	var inp structs.ReportInput
	json.NewDecoder(r.Body).Decode(&inp)

	security.CheckKey(append([]byte(``), inp.Key...), append([]byte(``), inp.Reporter...))

	c := make(chan string)

	// Unreport
	go func() {
		delReportUrl := "http://localhost:9200/asamahe/report/_query"
		delReportReq := delReportReqJsonTpl(inp.Reporter, inp.JokeId)
		err := sunhttp.DeleteStruct(delReportUrl, delReportReq)
		if err != nil {
			log.Printf("\n%s", err)
		} else {
			fmt.Fprintf(w, "{\"Status\":0}")
		}
		c <- "unreport"
	}()

	// Reduce ReportCount
	go func() {
		url2 := ("http://localhost:9200/asamahe/joke/" + inp.JokeId + "/_update")
		addReportCountReq := []byte(`{"script" : "ctx._source.ReportCount-=1"}`)
		err := sunhttp.PostNoResponse(url2, addReportCountReq)
		if err != nil {
			log.Printf("\n%s", err)
		}
		c <- "reduce ReportCount"
	}()

	for i1 := 0; i1 < 2; i1++ {
		<-c
	}
}

func delReportReqJsonTpl(reporter, jokeId string) []byte {
	b := []byte(`{"query":{"bool":{"must":[{"match":{"Reporter":"`)
	b = append(b, reporter...)
	b = append(b, []byte(`"}},{"match":{"JokeId":"`)...)
	b = append(b, jokeId...)
	b = append(b, []byte(`"}}]}}}`)...)
	return b
}
