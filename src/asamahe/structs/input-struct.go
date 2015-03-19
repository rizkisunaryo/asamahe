package structs

type SetUserInput struct {
	Id     string
	Key    string
	Name   string
	PicUrl string
}

type JokeInput struct {
	Key    string
	Viewer string
	JokeId string
	Lat    float64
	Long   float64
}

type NewJokesInput struct {
	Key      string
	Viewer   string
	Time     string
	IsBefore int
	Lat      float64
	Long     float64
}

type HotJokesInput struct {
	Key      string
	Viewer   string
	Score    int
	IsBefore int
	From     int
	Size     int
	Lat      float64
	Long     float64
}

type CreateJokeInput struct {
	Key     string
	Joker   string
	Title   string
	Content string
	Lat     float64
	Long    float64
}

type DeleteJokeInput struct {
	Key    string
	UserId string
	JokeId string
	Lat    float64
	Long   float64
}

type CommentInput struct {
	Key         string
	JokeId      string
	Commentator string
	Comment     string
	Lat         float64
	Long        float64
}

type DeleteCommentInput struct {
	Key       string
	UserId    string
	CommentId string
	JokeId    string
	Lat       float64
	Long      float64
}

type LikeInput struct {
	Key    string
	JokeId string
	Liker  string
	Lat    float64
	Long   float64
}

//type DislikeInput struct {
//	Key    string
//	UserId string
//	LikeId string
//	JokeId string
//	Lat    float64
//	Long   float64
//}

type ReportInput struct {
	Key      string
	JokeId   string
	Reporter string
	Lat      float64
	Long     float64
}

//type UnreportInput struct {
//	Key      string
//	UserId   string
//	ReportId string
//	JokeId   string
//	Lat      float64
//	Long     float64
//}

type SearchInput struct {
	Key      string
	Viewer   string
	Keywords []string
	Lat      float64
	Long     float64
}

type UserJokesInput struct {
	Key      string
	Viewer   string
	Joker    string
	Time     string
	IsBefore int
	Lat      float64
	Long     float64
}

type TopJokersInput struct {
	Key    string
	Viewer string
	Lat    float64
	Long   float64
}
