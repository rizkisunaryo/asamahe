package structs

type JokeOutput struct {
	Status   int
	Joke     Joke
	Comments []Comment
	Likes    []Like
	Reports  []Report
}

type JokesOutput struct {
	Jokes []Joke
}

type JokerJokesOutput struct {
	Joker User
	Jokes []Joke
}

type JokersOutput struct {
	Jokers []User
}
