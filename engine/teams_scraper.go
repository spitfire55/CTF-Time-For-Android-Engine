package engine

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"fmt"
)

type Team struct { // key = Team ID
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

/*
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
	Url   string
	Image string // relative Url to image
}

type Writeup struct { // key = Writeup ID
	TaskId     string
	AuthorTeam int //Team Id of author
	WriteupId  int
	WriteupUrl string
}

type Task struct { // Key = Task ID
	EventId     int // Event Id
	Tags        []string
	Description string
	Points      int // TODO: Some of the point values are # + # instead of just #
}
*/

func ParseAndStoreTeam(resp *http.Response) error {
	var team Team
	rootSel, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	team.Name = rootSel.Find("ul.breadcrumb li.active").Text()

	rootSel.Find(".span10 h5").Siblings().Find("li").Each(func(i int, selection *goquery.Selection) {
		team.Aliases = append(team.Aliases, selection.Text())
	})

	rootSel.Find(".span10 p").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()
		if selection.Parent().Is(".well") {
			team.Description = text
		} else if strings.Contains(text, "Academic team ") {
			team.Academic = strings.Join(strings.Split(text, " ")[2:], " ")
		} else if strings.Contains(text, "Website: ") && team.Website == "" {
			team.Website = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Twitter: ") {
			team.Twitter = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Skype: ") {
			team.Skype = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Telegram: ") {
			team.Telegram = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "ICQ: ") {
			team.ICQ = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Email: ") {
			team.Email = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "LinkedIn: ") {
			team.LinkedIn = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Jabber: ") {
			team.Jabber = strings.Join(strings.Split(text, " ")[1:], " ")
		} else if strings.Contains(text, "Other: ") {
			team.OtherLinks = append(team.OtherLinks, strings.Join(strings.Split(text, " ")[1:], " "))
		}
	})

	team.Logo, _ = rootSel.Find(".span2 img").First().Attr("src")
	return nil
}
