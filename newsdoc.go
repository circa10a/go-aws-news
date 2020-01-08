package awsnews

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func getNewsDoc(year int, month int) *goquery.Document {
	url := fmt.Sprintf("https://aws.amazon.com/about-aws/whats-new/%v/%02d/", year, month)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func (d *newsDoc) getSelectionItems() *goquery.Selection {
	return d.Find("li.whats-new")
}

func (d *newsDoc) getSelectionTitle(s *goquery.Selection) string {
	return s.Find("h3").Text()
}

func (d *newsDoc) getSelectionItemLink(s *goquery.Selection) string {
	link, found := s.Find("a").Attr("href")
	if !found {
		link = "Not found"
	}
	return link
}

func (d *newsDoc) getSelectionItemDate(s *goquery.Selection) string {
	return s.Find("div.date").Text()
}
