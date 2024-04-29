package main

import (
	"log"
	"net/http"

	pack "github.com/Yeosu-expo/highNoonServer/packages"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/ServingChunk", func(w http.ResponseWriter, r *http.Request) {
		pack.ServingChunkHandler(w, r)
	}).Methods("POST")

	router.HandleFunc("/sign_up", func(w http.ResponseWriter, r *http.Request) {
		pack.SignUpHandler(w, r)
	}).Methods("POST")

	router.HandleFunc("/sign_in", func(w http.ResponseWriter, r *http.Request) {
		pack.SignInHandler(w, r)
	}).Methods("POST")

	if err := http.ListenAndServe("210.125.31.150:6000", router); err != nil {
		log.Println("failed to open server:", err)
		return
	}
}
