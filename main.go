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
	}).Methods("POST")

	router.HandleFunc("/sign_up", func(w http.ResponseWriter, r *http.Request) {
		pack.SignUpHandler(w, r)
	}).Methods("POST")

	router.HandleFunc("/sign_in", func(w http.ResponseWriter, r *http.Request) {
		pack.SignInHandler(w, r)
	}).Methods("POST")

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

	router.HandleFunc("/insertChost", pack.InsertChostHandler).Methods("POST")
	router.HandleFunc("/getChost", pack.GetChostHandler).Methods("GET")
	router.HandleFunc("/getChost2", pack.GetChostRegardAccuracyHandler).Methods("GET")

	if err := http.ListenAndServe(":6000", router); err != nil {
		log.Println("failed to open server:", err)
		return
	}
}
