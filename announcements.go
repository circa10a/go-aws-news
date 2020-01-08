package awsnews

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Fetch Get an array of Announcements
func Fetch(year int, month int) Announcements {
	var Announcements []announcement
	doc := newsDoc{getNewsDoc(year, month)}
	doc.getSelectionItems().Each(func(i int, s *goquery.Selection) {
		title := doc.getSelectionTitle(s)
		link := fmt.Sprintf("https:%v", doc.getSelectionItemLink(s))
		date := strings.TrimSpace(doc.getSelectionItemDate(s))
		Announcements = append(Announcements, announcement{Title: title, Link: link, PostDate: date})
	})
	return Announcements
}
