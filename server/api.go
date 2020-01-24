package main

import (
	"github.com/JokeTrue/Devim-Test-Case/shared"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func initAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/api/numbers/list/", getsavedNumbers).Methods("GET")

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Printf("ERROR: Server init failed, %s", err.Error())
		os.Exit(1)
	}
}

func getsavedNumbers(w http.ResponseWriter, r *http.Request) {
	shared.Respond(w, http.StatusOK, savedNumbers)
}
