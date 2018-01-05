package engine

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/html"
)

// A Ranking contains the fields to be stored in Firestore for a
type Ranking struct {
	Rank        int
	TeamName    string
	TeamUrl     string
	CountryFlag string
	CountryID   string
	Score       float64
}

// ParseAndStoreRankings parses the response body of a rankings page for the values needed to create a Ranking struct. Once
// the page is parsed and an array of Rankings is created, the results are stored in the Firestore database IF the sha256
// hash of Rankings array is different from the hash stored in the Firestore database. If the page is new, no hash exists
// in the Firestore database, which means a page hash document will be created along with the Ranking documents.
func ParseAndStoreRankings(response *http.Response, pageNumber int, year string, fbc FirebaseContext) error {
	//TODO: Migrate from html.Tokenizer to GoQuery
	var rankings []Ranking
	pageNumDoc := fmt.Sprintf("Page%dHash", pageNumber)
	firstRow := true // flag to ignore table headers
	z := html.NewTokenizer(response.Body)

	for {
		tk := z.Next()
		switch {
		case tk == html.ErrorToken: // reached end of HTML file
			goto hashCheck
		case tk == html.StartTagToken:
			if string(z.Raw()) == "<tr>" {
				var ranking Ranking
				var rankingRow []html.Token
				rowLength := 0

				if firstRow {
					firstRow = false
					break
				}
				z.Next()

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

hashCheck:
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
