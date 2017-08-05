package main

import (
	"fmt"
	"net/http"
	//"github.com/anaskhan96/soup"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"log"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := &http.Client {
		Transport: &urlfetch.Transport{
			Context: ctx,
			// local dev app server doesn't like Lets Encrypt certs...
			AllowInvalidServerCertificate: appengine.IsDevAppServer(),
		},
	}
	resp, err := client.Get("https://ctftime.org/api/v1/top/",)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc := string(body)
	fmt.Fprint(w, doc)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8080", nil)
	appengine.Main()
}
