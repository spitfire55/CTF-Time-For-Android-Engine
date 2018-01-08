package main

import (
	"CTF-Time-For-Android-Engine/engine"
	"net/http"
)

// defaultHandler handles any requests to the engine that are not defined by other handlers. By default, any undefined
// request does nothing. This default handler should never be changed.
func defaultHandler(_ http.ResponseWriter, _ *http.Request) {}

func main() {
	http.HandleFunc("/ctfs", engine.UpdateCtfsHandler)
	http.HandleFunc("/rankings", engine.UpdateRankingsHandler)
	http.HandleFunc("/teams", engine.UpdateTeamsHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe("localhost:8080", nil)
}
