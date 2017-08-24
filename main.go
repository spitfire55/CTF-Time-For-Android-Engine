/*
This project enables storing data from ctftime.org in a Firebase database.
For some aspects of the website, the API does not easily expose the data needed to build
a robust mobile application. As a result, some web scraping is done in order to
acquire content such as writeups.

The following license applies to all main package Go files:

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

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/current-rankings", checkCurrentRankingsHandler)
}

func setup(url string, ctx context.Context, w http.ResponseWriter) []byte {
	client := &http.Client {
		Transport: &urlfetch.Transport{
			Context: ctx,
			AllowInvalidServerCertificate: appengine.IsDevAppServer(),
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	body := setup("https://ctftime.org/api/v1/top/", ctx, w)
	ranking := getAllRankings(body)
	connect(ctx)
	saveAllRankings(ranking, ctx)
}

func checkCurrentRankingsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	body := setup("https://ctftime.org/api/v1/top/2017/", ctx, w)
	ranking := getCurrentRankings(body)
	connect(ctx)
	saveCurrentRankings(ranking, ctx)
}

func main() {
	appengine.Main()
}
