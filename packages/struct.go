package packages

type PlayerChunk struct {
	Id       string `json:"id"`
	DateTime string `json:"dateTime"`
	Score    int    `json:"score"`
}

type UserInfoChunk struct {
	User_ID       string `json:"user_id"`
	User_Password string `json:"user_password"`
}

type MatchInfo struct {
	uid     string
	mmr     int
	resChan *chan MatchInfo
}

type MatchRes struct {
	Rival string `json:"rival"`
}

type RankInfo struct {
	UID string `json:"uid"`
	MMR int    `json:"mmr"`
}
