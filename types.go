package awsnews

import "github.com/PuerkitoBio/goquery"

// Announcements Array of the type []announcement which contains
// an announcement title, link and postdate
type Announcements []announcement

// announcement Hold significant data for each aws announcement
type announcement struct {
	Title    string
	Link     string
	PostDate string
}

// newDoc Hold html document
type newsDoc struct {
	*goquery.Document
}
