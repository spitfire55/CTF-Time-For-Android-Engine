package engine

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// A Ranking represents a row in the ctftime.org Rankings
type Ranking struct {
	Rank      int
	TeamName  string
	TeamId    int
	CountryId string
	Score     float64
	Events    int
}

// ParseAndStoreRankings parses the response body of a rankings page for the values needed to create a Ranking struct. Once
// the page is parsed and an array of Rankings is created, the results are stored in the Firestore database IF the sha256
// hash of the Rankings array is different from the hash stored in the Firestore database for the particular page. If the
// final page parsed is a newly created page, no page hash document exists in the Firestore database. We then create a page
// hash document along with the Ranking documents.
// NOTE: The page hash is computed from an array of Ranking types (i.e. after parsing), not from the response body of the
// request.
func ParseAndStoreRankings(response *http.Response, pageNumber int, year string, fbc FirebaseContext) error {
	var rankings []Ranking
	pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
	rootSel, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return err
	}

	rootSel.Find(".table.table-striped tr").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			return
		}
		columns := selection.Find("td")
		rank, _ := strconv.Atoi(columns.Nodes[0].Data)
		url, _ := columns.Find("a").Attr("href")
		id, _ := strconv.Atoi(strings.Split(url, "/")[1])
		score, _ := strconv.ParseFloat(columns.Nodes[3].Data, 64)
		events := 0
		if columns.Length() == 5 {
			events, _ = strconv.Atoi(columns.Nodes[4].Data)
		}
		rankings = append(rankings, Ranking{
			Rank:      rank,
			TeamName:  columns.Nodes[1].Data,
			TeamId:    id,
			CountryId: columns.Find("img").AttrOr("alt", ""),
			Score:     score,
			Events:    events,
		})
	})

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
