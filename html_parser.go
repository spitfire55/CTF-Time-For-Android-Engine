package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/net/html"
	"net/http"
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

func parseAndStoreRankings(response *http.Response, year int, pageNumber int, fbc FirebaseContext) error {
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
				fmt.Printf("%#v\n", rankingRow)
				ranking.Rank, _ = strconv.Atoi(rankingRow[1].Data)
				ranking.TeamUrl = rankingRow[4].Attr[0].Val
				ranking.TeamName = rankingRow[5].Data
				if rowLength == 17 { // 2017 w/ flag
					ranking.CountryFlag = rankingRow[9].Attr[0].Val
					ranking.CountryID = rankingRow[9].Attr[1].Val
					ranking.Score, _ = strconv.ParseFloat(rankingRow[12].Data, 64)
				} else if year == 16 { // 2017 w/o flag
					ranking.Score, _ = strconv.ParseFloat(rankingRow[11].Data, 64)
				} else if rowLength == 15 { // < 2017 w/ flag
					ranking.CountryFlag = rankingRow[9].Attr[]
				}
				fmt.Printf("%#v\n", ranking)
				rankings = append(rankings, ranking)
			}
		}
	}
finish:
	if len(rankings) > 0 {
		sha256Hash := sha256.New()
		sha256Hash.Write([]byte(fmt.Sprintf("%#v", rankings)))
		resultsHash := base64.StdEncoding.EncodeToString(sha256Hash.Sum(nil))
		hashDiff, err := hashDiff(resultsHash, pageNumber, year, fbc)
		if err != nil && !hashDiff {
			return err
		}
		pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
		if hashDiff {
			storeRankingsHash(resultsHash, pageNumDoc, year, fbc)
			storeRankings(rankings, year, fbc)
		}
		return nil
	}
	return errors.New("length of rankings array is zero")
}
