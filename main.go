package main

import (
	"CTF-Time-For-Android-Engine/engine"
	"net/http"
)

func defaultHandler(_ http.ResponseWriter, _ *http.Request) {}

func main() {
	http.HandleFunc("/rankings", engine.UpdateCurrentRankingsHandler)
	http.HandleFunc("/teams", engine.UpdateTeamsHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe("localhost:8080", nil)
}
