package structs

type UserResponse struct {
	Source User `json:"_source"`
}

type CreateJokeResponse struct {
	Id string `json:"_id"`
}

type JokeQueryResponse struct {
	Hits JokeHitsChildren `json:"hits"`
}

type JokeHitsChildren struct {
	Hits []JokeResponse `json:"hits"`
}

type JokeResponse struct {
	Id     string `json:"_id"`
	Source Joke   `json:"_source"`
}

type IdResponse struct {
	Id string `json:"_id"`
}

type CommentQueryResponse struct {
	Hits CommentHitsChildren `json:"hits"`
}

type CommentHitsChildren struct {
	Hits []CommentResponse `json:"hits"`
}

type CommentResponse struct {
	Id     string  `json:"_id"`
	Source Comment `json:"_source"`
}

type CheckResponse struct {
	Hits Hits `json:"hits"`
}

type LikeQueryResponse struct {
	Hits LikeHitsChildren `json:"hits"`
}

type LikeHitsChildren struct {
	Hits []LikeResponse `json:"hits"`
}

type LikeResponse struct {
	Source Like `json:"_source"`
}

type ReportQueryResponse struct {
	Hits ReportHitsChildren `json:"hits"`
}

type ReportHitsChildren struct {
	Hits []ReportResponse `json:"hits"`
}

type ReportResponse struct {
	Source Report `json:"_source"`
}

type TopJokersResponse struct {
	Aggregations TopJokersAggregations `json:"aggregations"`
}

type TopJokersAggregations struct {
	Group TopJokersGroup `json:"group_by_joker"`
}

type TopJokersGroup struct {
	JokerScores []JokerScore `json:"buckets"`
}
