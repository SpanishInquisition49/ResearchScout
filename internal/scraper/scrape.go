package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unipi-research-crawler/internal/database"

	"github.com/gocolly/colly"
)

const (
	baseUrl = "https://bandi.unipi.it"
)

var db *database.Queries = nil

func Scrape(database *database.Queries) {
	c := colly.NewCollector()

	db = database
	c.OnHTML("article", getCallDetails)
	c.OnHTML("tbody", getCalls)
	err := c.Visit(fmt.Sprintf("%s/public/Bandi?type=BRS&str=34&s=0", baseUrl))
	if err != nil {
		log.Fatalf("Failed to visit the page: %v", err)
	}
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
	call := database.CreateCallParams{
		Title:        title,
		Deadline:     expirationDate,
		Requirements: callUrl,
		ApplyModule:  applicationUrl,
	}

	// insert the call into the database
	if db == nil {
		log.Fatal("Database is not initialized")
	}
	ctx := context.Background()
	if _, err := db.CreateCall(ctx, call); err != nil {
		log.Printf("Failed to create call: %v", err)
	}
}
