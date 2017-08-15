/*
Package ctftime enables storing data from ctftime.org in a Firebase database.
For some aspects of the website, the API does not easily expose the data needed to build
a robust mobile application. As a result, some web scraping is done in order to
acquire content such as writeups.

The following license applies to all package files:

Copyright 2017 Dale Lakes

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of
the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/
package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
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
	ranking := getTeamRankings(body)
	connect(ctx)
	saveTeams(ranking)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe("localhost:8080", nil)
	appengine.Main()
}
