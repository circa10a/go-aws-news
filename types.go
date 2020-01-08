package awsnews

import "github.com/PuerkitoBio/goquery"

//News Create instance to fetch announcements
type News struct {
	Year  int
	Month int
}

// Announcements Array of the type []announcement which contains
// an announcement title, link and postdate
type Announcements []announcement

type announcement struct {
	Title    string
	Link     string
	PostDate string
}

type newsDoc struct {
	*goquery.Document
}
