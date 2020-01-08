package awsnews

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Get Get an array of Announcements
func (n *News) Get() Announcements {
	var Announcements []announcement

	doc := newsDoc{getNewsDoc(n.Year, n.Month)}
	doc.getSelectionItems().Each(func(i int, s *goquery.Selection) {
		title := doc.getSelectionTitle(s)
		link := fmt.Sprintf("https:%v", doc.getSelectionItemLink(s))
		date := strings.TrimSpace(doc.getSelectionItemDate(s))
		Announcements = append(Announcements, announcement{Title: title, Link: link, PostDate: date})
	})
	return Announcements
}

// Print Print all news results
func (n *News) Print() {
	for _, v := range n.Get() {
		fmt.Printf("Title: %v\n", v.Title)
		fmt.Printf("Link: %v\n", v.Link)
		fmt.Printf("Date: %v\n", v.PostDate)
		fmt.Println()
	}
}
