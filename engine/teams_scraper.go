package engine

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type Team struct { // key = Team ID
	hash        string
	Name        string
	Aliases     []string
	Academic    string
	Description string
	Logo        string // relative URL
	Website     string
	Twitter     string
	Email       string
	ICQ         string
	Skype       string
	LinkedIn    string
	Telegram    string
	Jabber      string
	OtherLinks  []string
	Member      []string
}

func ParseAndStoreTeam(teamId int, resp *http.Response, fbc FirebaseContext) error {
	var team Team
	rootSel, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	team.Name = rootSel.Find("h2").Text()

	rootSel.Find(".span10 h5").Siblings().Find("li").Each(func(i int, selection *goquery.Selection) {
		team.Aliases = append(team.Aliases, selection.Text())
	})

	rootSel.Find(".span10 p").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()
		if selection.Parent().Is(".well") {
			team.Description = text
		} else if strings.Contains(text, "Academic team ") {
			team.Academic = strings.Join(strings.Split(text, " ")[2:], " ")
		} else if strings.Contains(text, "Website: ") && team.Website == "" { // keep oldest
			team.Website = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Twitter: ") {
			team.Twitter = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Skype: ") {
			team.Skype = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Telegram: ") {
			team.Telegram = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "ICQ: ") {
			team.ICQ = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Email: ") && team.Email == "" { // keep oldest
			team.Email = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "LinkedIn: ") {
			team.LinkedIn = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Jabber: ") {
			team.Jabber = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Other: ") { // keep all of them
			team.OtherLinks = append(team.OtherLinks, strings.Join(strings.Split(text, " ")[1:], " "))
		}
	})

	team.Logo, _ = rootSel.Find(".span2 img").First().Attr("src")

	teamHash := CalculateHash(team)
	team.hash = teamHash
	hashDiff, err := TeamHashDiff(teamId, team, fbc)
	if err != nil {
		return err
	}
	if hashDiff {
		StoreTeam(teamId, team, fbc)
	}
	return nil
}
