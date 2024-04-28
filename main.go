package main

import (
	"log"
	"net/http"

	pack "github.com/Yeosu-expo/highNoonServer/packages"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	tmp := make(chan pack.MatchInfo)
	matchingChan := &tmp

	go pack.Matching(matchingChan)

	router.HandleFunc("/ServingChunk", func(w http.ResponseWriter, r *http.Request) {
		pack.ServingChunkHandler(w, r)
	})

	router.HandleFunc("/playResult", func(w http.ResponseWriter, r *http.Request) {
		pack.PlayResultHandler(w, r)
	}).Methods("GET")

	router.HandleFunc("/realTimeMatching", func(w http.ResponseWriter, r *http.Request) {
		pack.RealTimeMatchingHandler(w, r, matchingChan)
	})

	router.HandleFunc("/matching", func(w http.ResponseWriter, r *http.Request) {
		pack.MatchingHandler(w, r)
	})

	router.HandleFunc("/getRank", func(w http.ResponseWriter, r *http.Request) {
		pack.GetRankHandler(w, r)
	})

	if err := http.ListenAndServe(":6000", router); err != nil {
		log.Println("failed to open server:", err)
		return
	}
}
