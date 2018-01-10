// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Ranking struct {
	CountryId string
	Events    int
	Rank      int
	Score     float64
	TeamId    int
	TeamName  string
}

func ParseAndStoreRankings(response *http.Response, pageNumber int, year string, fbc FirebaseContext) error {
	var rankings []Ranking
	pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
	rootSel, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return err
	}

	rootSel.Find(".table.table-striped tr").Each(func(rowIndex int, rowSelection *goquery.Selection) {
		var ranking Ranking
		if rowIndex == 0 {
			return
		}
		rowSelection.Find("td").Each(func(colIndex int, colSelection *goquery.Selection) {
			switch colIndex {
			case 0:
				ranking.Rank, _ = strconv.Atoi(colSelection.Text())
			case 1:
				ranking.TeamName = colSelection.Find("a").Text()
				teamIdUrl, _ := colSelection.Find("a").Attr("href")
				teamIdSplit := strings.Split(teamIdUrl, "/")
				ranking.TeamId, _ = strconv.Atoi(teamIdSplit[len(teamIdSplit)-1])
			case 2:
				ranking.CountryId, _ = colSelection.Find("img").Attr("alt")
			case 3:
				ranking.Score, _ = strconv.ParseFloat(colSelection.Text(), 64)
			case 4:
				ranking.Events, _ = strconv.Atoi(colSelection.Text())
			}
		})
		rankings = append(rankings, ranking)
	})
	fmt.Printf("%#v\n", rankings)

	if len(rankings) > 0 {
		resultsHash := CalculateHash(rankings)
		hashDiff, err := RankingsHashDiff(resultsHash, pageNumber, year, fbc)
		if err != nil {
			return err
		}
		if hashDiff {
			StoreRankingsHash(resultsHash, pageNumDoc, year, fbc)
			StoreRankings(rankings, year, fbc)
		}
	} else {
		return errors.New("length of rankings array is zero")
	}
	return nil
}
