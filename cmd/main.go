package main

import (
	"log"

	scrapeOMXHelsinki "github.com/mule/go-scrape-omx-helsinki"
)

func main() {

	log.Println("Starting OMX Helsinki Scraper")
	log.Println(scrapeOMXHelsinki.Config())

}
