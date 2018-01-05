package engine

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type Team struct { // key = Team ID
	Name        string
	Aliases     []string
	Academic    string
	Description string
	Logo        string // relative URL
	Website     string
	Twitter     string
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
	//TODO: Migrate from soup to GoQuery
	var team Team
	rootSel, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	teamHeader := rootSel.Find("h2").First()
	if teamHeader.Find("img").Length() != 0 {
		team.Name = string([]rune(teamHeader.Text())[1:])
	} else {
		team.Name = teamHeader.Text()
	}
	fmt.Println(team.Name)
	return nil
	/*
			// Find Aliases, Academic, Website, and Twitter
			spanTens := root.FindAll("div", "class", "span10")
			fmt.Printf("%#v", spanTens)
			for _, spanTen := range spanTens {
				aliases := findAliases(spanTen)
				if len(aliases) != 0 {
					team.Aliases = aliases
					break
				}
				academic := findAcademic(spanTen)
				if academic != "" {
					team.Academic = academic
					break
				}
				website, twitter := findWebStuff(spanTen)
				if website != "" {
					team.Website = website
				}
				if twitter != "" {
					team.Twitter = twitter
				}
			}

			// Find Description
			description := root.Find("div", "class", "well").Find("p").Text()
			if description != "" {
				team.Description = description
			}

			// Find Logo (default or non-default)
			team.Logo = root.Find("div", "class", "span2").Find("img").Attrs()["src"]

			// Find Team Members
			teamTable := root.Find("div", "id", "recent_members").FindAll("a")
			for _, teamMember := range teamTable {
				team.Member = append(team.Member, teamMember.Text())
			}
			results
		}

		func findAliases(dom soup.Root) []string {
			var aliases []string
			if dom.Find("h5").Text() == "Also known as" {
				domAliases := dom.FindAll("li")
				for _, domAlias := range domAliases {
					aliases = append(aliases, domAlias.Text())
				}
			}
			return aliases
		}

		func findAcademic(dom soup.Root) string {
			academic := ""
			if dom.Find("p").Find("b").Text() == "Academic team" {
				academic = dom.Find("p").Text()
			}
			return academic
		}

		func findWebStuff(dom soup.Root) (string, string) {
			website := ""
			twitter := ""
			entries := dom.FindAll("p")
			for _, entry := range entries {
				if entry.Text() == "Website: " && website != "" { // use first website only
					website = entry.Find("a").Attrs()["href"]
				}
				if entry.Text() == "Twitter: " {
					twitter = entry.Find("a").Attrs()["href"]
				}
			}
			return website, twitter
	*/
}
