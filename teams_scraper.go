// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

type Team struct { // key = Team ID
	Hash string
	// General
	Aliases             []string
	Academic            string
	CountryCode         string
	Description         string
	Logo                string // relative URL
	Members             []Member
	Name                string
	NameCaseInsensitive string
	Scores              map[string]Score
	// Social
	Email      string
	ICQ        string
	Jabber     string
	LinkedIn   string
	OtherLinks []string
	Skype      string
	Telegram   string
	Twitter    string
	Website    string
}

type Score struct {
	Points float64
	Rank   int
}

type Member struct {
	Id   int
	Name string
}

func ParseAndStoreTeam(teamId int, resp *http.Response, fbc FirebaseContext) error {
	var team Team
	team.Scores = make(map[string]Score)
	rootSel, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	rootSel.Find(".span10 h5").Siblings().Find("li").Each(func(i int, selection *goquery.Selection) {
		team.Aliases = append(team.Aliases, selection.Text())
	})

	team.CountryCode, _ = rootSel.Find("h2 img").Attr("alt")

	rootSel.Find(".span10 p").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()
		if selection.Parent().Is(".well") {
			team.Description = text
		} else if strings.Contains(text, "Academic team ") {
			team.Academic = strings.Join(strings.Split(text, " ")[2:], " ")
		} else if strings.Contains(text, "Email: ") && team.Email == "" { // keep oldest
			team.Email = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "ICQ: ") {
			team.ICQ = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Jabber: ") {
			team.Jabber = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "LinkedIn: ") {
			team.LinkedIn = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Other: ") { // keep all of them
			team.OtherLinks = append(team.OtherLinks, strings.Join(strings.Split(text, " ")[1:], " "))
		} else if strings.Contains(text, "Skype: ") {
			team.Skype = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Telegram: ") {
			team.Telegram = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Twitter: ") {
			team.Twitter = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Website: ") && team.Website == "" { // keep oldest
			team.Website = strings.Join(strings.Split(text, " ")[1:], " ")
		}
	})

	team.Logo, _ = rootSel.Find(".span2 img").First().Attr("src")

	rootSel.Find("#recent_members td").Each(func(i int, selection *goquery.Selection) {
		idUrl, _ := selection.Find("a").Attr("href")
		idSplit := strings.Split(idUrl, "/")
		id, _ := strconv.Atoi(idSplit[len(idSplit)-1])
		team.Members = append(team.Members, Member{Id: id, Name: selection.Find("a").Text()})
	})

	team.Name = rootSel.Find(".breadcrumb .active").Text()
	team.NameCaseInsensitive = strings.ToLower(team.Name)

	for year := 2011; year < 2018; year++ {
		yearStr := strconv.Itoa(year)
		findStr := fmt.Sprintf("#rating_%s p b", yearStr)
		yearPoints, _ := strconv.ParseFloat(rootSel.Find(findStr).Last().Text(), 64)
		yearRankStripped := strings.TrimSpace(rootSel.Find(findStr).First().Text())
		yearRank, _ := strconv.Atoi(yearRankStripped)
		team.Scores[yearStr] = Score{
			yearPoints,
			yearRank,
		}
	}

	teamHash := CalculateHash(team)
	team.Hash = teamHash
	hashDiff, err := CompareTeamHash(teamId, team, fbc)
	if err != nil {
		return err
	}
	if hashDiff {
		err = StoreTeam(teamId, team, fbc)
		if err != nil {
			return err
		}
	}
	return nil
}
