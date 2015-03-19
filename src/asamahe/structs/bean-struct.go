package structs

type Hits struct {
	Total int `json:"total"`
}

type User struct {
	Id     string
	Name   string
	PicUrl string
}

type Joke struct {
	JokeId       string
	Joker        string
	JokerName    string
	JokerPicUrl  string
	Title        string
	Content      string
	CommentCount int
	LikeCount    int
	ReportCount  int
	IsLiked      int
	IsReported   int
	Score        int
	IsAd         int
	IsBlocked    int
	Time         string
}

type Comment struct {
	JokeId            string
	CommentId         string
	Commentator       string
	CommentatorName   string
	CommentatorPicUrl string
	Comment           string
	IsBlocked         int
	Time              string
}

type Like struct {
	JokeId    string
	Liker     string
	LikerName string
	Time      string
}

type Report struct {
	JokeId       string
	Reporter     string
	ReporterName string
	Time         string
}

type JokerScore struct {
	Key      string `json:"key"`
	SumScore Score  `json:"sumScore"`
}

type Score struct {
	Value float64 `json:"value"`
}
