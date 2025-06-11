package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Call struct {
	Title          string
	ExpirationDate string
	CallUrl        string
	ModuleUrl      string
}

const (
	baseUrl = "https://bandi.unipi.it"
)

var calls []Call

func Scrape() *[]Call {
	calls = []Call{} // Reset calls slice to avoid duplicates
	c := colly.NewCollector()

	c.OnHTML("article", getCallDetails)
	c.OnHTML("tbody", getCalls)
	err := c.Visit(fmt.Sprintf("%s/public/Bandi?type=BRS&str=34&s=0", baseUrl))
	if err != nil {
		log.Fatalf("Failed to visit the page: %v", err)
	}

	// Create a new slice to give back to the caller
	res := make([]Call, len(calls))
	for i, call := range calls {
		res[i] = Call{
			call.Title,
			call.ExpirationDate,
			call.CallUrl,
			call.ModuleUrl,
		}
	}
	return &res
}

func getCalls(h *colly.HTMLElement) {
	// Print the element class
	h.ForEach("tr", getCall)
}

func getCall(_ int, h *colly.HTMLElement) {
	url := h.ChildAttr("a", "href")
	name := h.ChildText("a")

	if url == "" && name == "" {
		return
	}
	h.Request.Visit(fmt.Sprintf("%s%s", baseUrl, url))
}

func getCallDetails(h *colly.HTMLElement) {
	// check if the article contains 2 div with class "row mt-4"
	// if the length of childDivs is less than 2, return (we are not interested in this article)
	childDivs := h.ChildAttrs("div.row.mt-4", "class")
	if childDivs == nil || len(childDivs) < 2 {
		return
	}
	// the first div contains the title and the second div contains the call details
	titleFull := h.ChildText("div.mt-4.row div.col-12")
	expirationDate := h.ChildText("main.col-md-12 > div:nth-child(3) > article:nth-child(1) > span:nth-child(4)")

	// the call details url is the first link in the second div
	callUrl := h.ChildAttr("div.mt-4.row a:first-child", "href")
	// the application module url is the second link in the second div
	applicationUrl := h.ChildAttr("p.download-list:nth-child(5) > a:nth-child(1)", "href")

	start := strings.Index(titleFull, "“")
	end := strings.Index(titleFull, "”")
	title := strings.TrimSpace(titleFull[start+1 : end])
	call := Call{
		Title:          title,
		CallUrl:        callUrl,
		ExpirationDate: expirationDate,
		ModuleUrl:      applicationUrl,
	}
	calls = append(calls, call)
}

func (c *Call) String() string {
	return fmt.Sprintf("Title: %s\nExpiration Date: %s\nRequirements URL: %s\nApplication Module URL: %s\n", c.Title, c.ExpirationDate, c.CallUrl, c.ModuleUrl)
}

func (c *Call) ToBotStringHTML() string {
	msgText := fmt.Sprintf(
		"🎓 <b>%s</b>\n\n"+
			"• <b>Deadline:</b> %s\n",
		c.Title,
		c.ExpirationDate,
	)
	return strings.ToValidUTF8(msgText, "")
}
