package main

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := &http.Client{
		Transport: &urlfetch.Transport{
			Context: ctx,
			// local dev app server doesn't like Lets Encrypt certs...
			AllowInvalidServerCertificate: appengine.IsDevAppServer(),
		},
	}
	resp, err := client.Get("https://ctftime.org/api/v1/top/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	ranking := decode(&body)
	fmt.Printf("%v", ranking)
	firebaseContents := test(ctx)
	fmt.Fprint(w, firebaseContents)


}

func main() {
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe("localhost:8080", nil)
	appengine.Main()
}
