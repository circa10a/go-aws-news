package news

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// newsDoc Extends the goquery.Document to create more receiver functions on the *goquery.Document type.
type newsDoc struct {
	*goquery.Document
}

// getNewsDocYear Fetches html from AWS depending on the year/month specified.
func getNewsDocYear(year int) (*goquery.Document, error) {
	url := fmt.Sprintf("https://aws.amazon.com/about-aws/whats-new/%d", year)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// getNewsDocMonth Fetches html from AWS depending on the year/month specified.
func getNewsDocMonth(year int, month int) (*goquery.Document, error) {
	url := fmt.Sprintf("https://aws.amazon.com/about-aws/whats-new/%v/%02d/", year, month)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// getSelectionItems Filters <li> of announcements from the html document.
func (d *newsDoc) getSelectionItems() *goquery.Selection {
	return d.Find("li.whats-new")
}

// getSelectionTitle Filters <h3> headlines from the aws html document.
func (d *newsDoc) getSelectionTitle(s *goquery.Selection) string {
	return s.Find("h3").Text()
}

// getSelectionItemLink Filters <a href>'s for a link which contains more details about each announcement.
func (d *newsDoc) getSelectionItemLink(s *goquery.Selection) string {
	link, found := s.Find("a").Attr("href")
	if !found {
		link = "Not found"
	}
	return link
}

// getSelectionItemData Filters the date posted for each announcement.
func (d *newsDoc) getSelectionItemDate(s *goquery.Selection) string {
	return s.Find("div.date").Text()
}
