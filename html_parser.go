package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/net/html"
	"net/http"
	"time"
	"google.golang.org/api/admin/directory/v1"
)

/*
 * RANKINGS
 */

type Ranking struct {
	Rank        int
	TeamName    string
	TeamUrl     string
	CountryFlag string
	CountryID   string
	Score       float64
}

type Team struct { // key = Team ID
	Aliases     []string
	Academic    string
	Description string
	Website     string
	Member []string
}

type Event struct { // key = Event ID
	Start         time.Time
	End           time.Time
	CtfId         int // ctf ID
	Format        string
	Website       string
	FutureWeight  float64
	CurrentWeight float64
	OrganizerId   int // Team ID of organizer
	Description   string
	Prizes        string
	Scoreboard    []Score
}

type Score struct { // no key, only used in Scoreboard array in Event
	TeamId       int
	CtfPoints    float64
	RatingPoints float64
}

type CTF struct { // key = CTF ID
	Url string
	Image string // relative Url to image
}

type Writeup struct { // key = Writeup ID
	TaskId string
	AuthorTeam int //Team Id of author
	WriteupId int
	WriteupUrl string
}

type Task struct { // Key = Task ID
	EventId int // Event Id
	Tags []string
	Description string
	Points int // TODO: Some of the point values are # + # instead of just #
}

func parseAndStoreTeams(response *http.Response, teamId int, fbc FirebaseContext) error {

}

func parseAndStoreRankings(response *http.Response, pageNumber int, year string, fbc FirebaseContext) error {
	var rankings []Ranking
	z := html.NewTokenizer(response.Body)
	firstRow := true
	for {
		tk := z.Next()
		switch {
		case tk == html.ErrorToken: // reached end of HTML file
			goto finish
		case tk == html.StartTagToken:
			if string(z.Raw()) == "<tr>" {
				if firstRow {
					firstRow = false
					break
				}
				z.Next()
				rowLength := 0
				var ranking Ranking
				var rankingRow []html.Token
				for string(z.Raw()) != "</tr>" {
					rankingRow = append(rankingRow, z.Token())
					rowLength++
					z.Next()
				}
				ranking.Rank, _ = strconv.Atoi(rankingRow[1].Data)
				ranking.TeamUrl = rankingRow[4].Attr[0].Val
				ranking.TeamName = rankingRow[5].Data
				if rowLength == 17 { // 2017 w/ flag
					ranking.CountryFlag = rankingRow[9].Attr[0].Val
					ranking.CountryID = rankingRow[9].Attr[1].Val
					ranking.Score, _ = strconv.ParseFloat(rankingRow[12].Data, 64)
				} else if rowLength == 16 { // 2017 w/o flag
					ranking.Score, _ = strconv.ParseFloat(rankingRow[11].Data, 64)
				} else if rowLength == 14 { // < 2017 w/ flag
					ranking.CountryFlag = rankingRow[9].Attr[0].Val
					ranking.CountryID = rankingRow[9].Attr[1].Val
					ranking.Score, _ = strconv.ParseFloat(rankingRow[12].Data, 64)
				} else { // < 2017 w/o flag (row length = 13)
					ranking.Score, _ = strconv.ParseFloat(rankingRow[11].Data, 64)
				}
				rankings = append(rankings, ranking)
			}
		}
	}
finish:
	if len(rankings) > 0 {
		sha256Hash := sha256.New()
		sha256Hash.Write([]byte(fmt.Sprintf("%#v", rankings)))
		resultsHash := base64.StdEncoding.EncodeToString(sha256Hash.Sum(nil))
		hashDiff, err := rankingsHashDiff(resultsHash, pageNumber, year, fbc)
		if err != nil && !hashDiff {
			return err
		}
		pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
		//if hashDiff {
		storeRankingsHash(resultsHash, pageNumDoc, year, fbc)
		storeRankings(rankings, year, fbc)
		//}
		return nil
	}
	return errors.New("length of rankings array is zero")
}
