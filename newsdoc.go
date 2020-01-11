package awsnews

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// getNewsDoc Fetches html from AWS depending on the year/month specified
func getNewsDoc(year int, month int) (*goquery.Document, error) {
	url := fmt.Sprintf("https://aws.amazon.com/about-aws/whats-new/%v/%02d/", year, month)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

// getSelectionItems Filters <li> of announcements from html doc
func (d *newsDoc) getSelectionItems() *goquery.Selection {
	return d.Find("li.whats-new")
}

// getSelectionTitle Filters <h3> headlines from html doc
func (d *newsDoc) getSelectionTitle(s *goquery.Selection) string {
	return s.Find("h3").Text()
}

// getSelectionItemLink Filters <a href>'s for a link which contains more details about each announcement
func (d *newsDoc) getSelectionItemLink(s *goquery.Selection) string {
	link, found := s.Find("a").Attr("href")
	if !found {
		link = "Not found"
	}
	return link
}

//getSelectionItemData Filters the date posted for each announment
func (d *newsDoc) getSelectionItemDate(s *goquery.Selection) string {
	return s.Find("div.date").Text()
}
