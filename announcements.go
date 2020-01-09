package awsnews

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
)

// Fetch Get an array of Announcements
func Fetch(year int, month int) Announcements {
	// Create []announcement
	var announcements Announcements
	doc := newsDoc{getNewsDoc(year, month)}
	doc.getSelectionItems().Each(func(i int, s *goquery.Selection) {
		title := doc.getSelectionTitle(s)
		link := fmt.Sprintf("https:%v", doc.getSelectionItemLink(s))
		date := parseDate(doc.getSelectionItemDate(s))
		announcements = append(announcements, announcement{Title: title, Link: link, PostDate: date})
	})
	return announcements
}

// Extract date from amazon's format
// Input: Posted on: Jan 7, 2020
// Output: Jan 7, 2020
func parseDate(postDate string) string {
	r, _ := regexp.Compile("[A-Z][a-z]{2}\\s[0-9]{1,2},\\s[0-9]{4}")
	return r.FindStringSubmatch(postDate)[0]
}

// ThisMonth get this month's AWS announcements
func ThisMonth() Announcements {
	currentTime := time.Now()
	return Fetch(currentTime.Year(), int(currentTime.Month()))
}

// Today get today's AWS announcments
func Today() Announcements {
	var todaysAnnouncements Announcements
	for _, announcement := range ThisMonth() {
		postDate, _ := time.Parse("Jan 2, 2006", announcement.PostDate)
		if dateEqual(postDate, time.Now()) {
			todaysAnnouncements = append(todaysAnnouncements, announcement)
		}
	}
	return todaysAnnouncements
}

// Print Print out ASCII table of your selection of AWS announcements
func (a *Announcements) Print() {
	data := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Title", "Date"})
	table.SetRowLine(true)
	for _, v := range *a {
		s := []string{
			v.Title,
			v.PostDate,
		}
		data = append(data, s)
	}
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
