package packages

type PlayerChunk struct {
	Id       string `json:"id"`
	DateTime string `json:"dateTime"`
	Score    int    `json:"score"`
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
