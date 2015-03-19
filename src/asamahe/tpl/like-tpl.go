package tpl

func LikeCheckReqJsonTpl(jokeId string, liker string) []byte {
	b := []byte(`{
  "query": {
    "bool": {
      "must": [
        { "match": { "JokeId": "`)
	b = append(b, jokeId...)
	b = append(b, []byte(`" } },
        { "match": { "Liker": "`)...)
	b = append(b, liker...)
	b = append(b, []byte(`" } }
      ]
    }
  }
}`)...)
	return b
}
