package tpl

func ReportCheckReqJsonTpl(jokeId string, reporter string) []byte {
	b := []byte(`{
  "query": {
    "bool": {
      "must": [
        { "match": { "JokeId": "`)
	b = append(b, jokeId...)
	b = append(b, []byte(`" } },
        { "match": { "Reporter": "`)...)
	b = append(b, reporter...)
	b = append(b, []byte(`" } }
      ]
    }
  }
}`)...)
	return b
}
