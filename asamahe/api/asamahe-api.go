package main

import (
	"asamahe/api/generic/comment"
	"asamahe/api/generic/joke"
	"asamahe/api/generic/like"
	"asamahe/api/generic/report"
	"asamahe/api/generic/search"
	"asamahe/api/generic/user"
	"net/http"
)

func main() {
	comment.HandleComment()
	comment.HandleDeleteComment()
	joke.HandleCreateJoke()
	joke.HandleDeleteJoke()
	joke.HandleHotJokes()
	joke.HandleJoke()
	joke.HandleNewJokes()
	like.HandleDislike()
	like.HandleLike()
	report.HandleReport()
	report.HandleUnreport()
	search.HandleSearch()
	user.HandleSetUser()
	user.HandleUserJokes()
	user.HandleTopJokers()

	err := http.ListenAndServe(":2903", nil)
	if err != nil {

	}
}
