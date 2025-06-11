package main

import (
	"log"
	"unipi-research-crawler/internal/scraper"
)

func main() {
	log.Println("Starting crawler...")
	calls := scraper.Scrape()
	for _, call := range *calls {
		log.Printf("-- Call Details --\n%s", call.String())
	}
	log.Println("Crawling completed.")
}
