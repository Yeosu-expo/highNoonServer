package packages

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defalutInc    = 10  // 기본적인 mmr 증감 값
	defalutIncDiv = 3   // mmr차이/deafalutIncDiv 를 증가함
	defalutDecDiv = 1   // mmr차이/deafalutDecDiv 를 감소함
	matchGap      = 100 // 매칭할 mmr 차이
)

// 결과에 따라 mmr을 올리고 내림
func PlayResultHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	rival := r.URL.Query().Get("rival")
	isWin := r.URL.Query().Get("isWin")

	db, err := sql.Open("mysql", "root:9250@tcp(localhost:3306)/highnoon")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	log.Println(uid, rival, isWin)
	res := getMMRInc(uid, rival, isWin, db)
	if res == 0 {
		return
	}

	var query string
	if isWin == "true" {
		query = "UPDATE MMR SET mmr=mmr+(?), winn=winn+1 WHERE uid=(?)"
	} else {
		query = "UPDATE MMR SET mmr=mmr-(?), losen=losen+1 WHERE uid=(?)"
	}
	abs := math.Abs(float64(res))
	rabs := int(abs)
	_, err = db.Exec(query, rabs, uid)
	if err != nil {
		log.Println("Failed to insert row:", err)
		return
	}

	log.Println(uid, "mmr +=", rabs)
}

// 매칭을 받고 고루틴으로 넘김
func RealTimeMatchingHandler(w http.ResponseWriter, r *http.Request, channel *chan MatchInfo) {
	uid := r.URL.Query().Get("uid")
	mmr := getMMR(uid)
	tmp := make(chan MatchInfo)

	info := MatchInfo{
		uid:     uid,
		mmr:     mmr,
		resChan: &tmp,
	}

	*channel <- info
	log.Println(uid, "(", mmr, ") ", "matching ~")

	rival := <-tmp
	log.Println(uid, "match with", rival.uid, rival.mmr)

	chunk := MatchRes{
		Rival: rival.uid,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(chunk)
	if err != nil {
		log.Println(err)
		return
	}
}

func MatchingHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	mmr := getMMR(uid)

	db, err := sql.Open("mysql", "root:9250@tcp(localhost:3306)/highnoon")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	query := "SELECT uid, mmr FROM mmr WHERE mmr >= (?) AND mmr <= (?) AND uid != (?)"
	rows, err := db.Query(query, mmr-matchGap, mmr+matchGap, uid)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	rows.Next()
	var rival string
	var rmmr int
	if err = rows.Scan(&rival, &rmmr); err != nil {
		log.Println(err)
		return
	}

	chunk := &RankInfo{
		UID: rival,
		MMR: rmmr,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(chunk)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%s(%d) match with %s(%d).\n", uid, mmr, rival, rmmr)
}

// mmr기반 rank정보를 출력 함
// 혹은 승수나 일정 기준으로 바꾸는 게 좋아보임
func GetRankHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:9250@tcp(localhost:3306)/highnoon")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	query := "SELECT uid, mmr from mmr ORDER BY mmr DESC"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	ranks := make([]RankInfo, 0)
	for rows.Next() {
		var uid string
		var mmr int
		rows.Scan(&uid, &mmr)

		info := RankInfo{
			UID: uid,
			MMR: mmr,
		}

		ranks = append(ranks, info)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ranks)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(ranks)
}

// 고루틴으로 실행되며
// 적절한 mmr의 범주의 플레이어가 매칭에 잡힐 때 결과를 각 플레이어에게 반환 함
func Matching(channel *chan MatchInfo) {
	players := make([]MatchInfo, 0)
	for {
		newInfo := <-*channel
		isMatched := false
		for _, info := range players {
			tmp := info.mmr - newInfo.mmr
			tmp2 := math.Abs(float64(tmp))
			gap := int(tmp2)

			if gap <= matchGap {
				*newInfo.resChan <- info
				*info.resChan <- newInfo

				isMatched = true
				break
			}
		}

		if !isMatched {
			players = append(players, newInfo)
		}
	}
}

func getMMR(uid string) int {
	db, err := sql.Open("mysql", "root:9250@tcp(localhost:3306)/highnoon")
	if err != nil {
		log.Println(err)
		return 0
	}
	defer db.Close()

	query := "SELECT mmr FROM mmr WHERE uid=(?)"
	rows, err := db.Query(query, uid)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer rows.Close()

	var ummr int
	rows.Next()
	if err = rows.Scan(&ummr); err != nil {
		log.Println(err)
		return 0
	}

	return ummr
}

func getMMRInc(uid string, rival string, isWin string, db *sql.DB) int {
	res := 0

	query := "SELECT mmr FROM mmr WHERE uid=(?) || uid=(?)"
	rows, err := db.Query(query, uid, rival)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer rows.Close()

	var ummr, rmmr int
	rows.Next()
	if err = rows.Scan(&ummr); err != nil {
		log.Println(err)
		return 0
	}
	rows.Next()
	if err = rows.Scan(&rmmr); err != nil {
		log.Println(err)
		return 0
	}

	tmp := math.Abs(float64(ummr - rmmr))
	mmrDiff := int(tmp)
	if isWin == "true" {
		res = defalutInc + mmrDiff/defalutIncDiv
	} else {
		res = defalutInc + mmrDiff/defalutDecDiv
	}

	return res
}
