package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"io/ioutil"
	"net/http"
)

func recursiveTeamCheck(ctx context.Context, w http.ResponseWriter) {
	highestNode := getHighestNode(ctx)
	body := checkNewTeam(ctx, highestNode)
	if body != nil {
		saveNewTeam(body, ctx)
		recursiveTeamCheck(ctx, w)
	}
}

func checkNewTeam(ctx context.Context, highestNode string) []byte {
	client := &http.Client{
		Transport: &urlfetch.Transport{
			Context: ctx,
			AllowInvalidServerCertificate: appengine.IsDevAppServer(),
		},
	}
	resp, err := client.Get("https://ctftime.org/api/v1/teams/" + highestNode + "/")
	if err != nil || resp.StatusCode != 200 {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body
}
