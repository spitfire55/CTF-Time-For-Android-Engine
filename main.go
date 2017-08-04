package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

)

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Get("https://www.google.com/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8080", nil)
	appengine.Main()
}
