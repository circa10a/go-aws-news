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
func Fetch(year int, month int) (Announcements, error) {
	// Create []announcement
	var announcements Announcements
	doc, err := getNewsDoc(year, month)
	if err != nil {
		return announcements, err
	}
	news := newsDoc{doc}
	news.getSelectionItems().Each(func(i int, s *goquery.Selection) {
		title := news.getSelectionTitle(s)
		link := fmt.Sprintf("https:%v", news.getSelectionItemLink(s))
		date := parseDate(news.getSelectionItemDate(s))
		announcements = append(announcements, announcement{Title: title, Link: link, PostDate: date})
	})
	return announcements, nil
}

// Extract date from amazon's format
// Input: Posted on: Jan 7, 2020
// Output: Jan 7, 2020
func parseDate(postDate string) string {
	r, _ := regexp.Compile("[A-Z][a-z]{2}\\s[0-9]{1,2},\\s[0-9]{4}")
	return r.FindStringSubmatch(postDate)[0]
}

// ThisMonth get this month's AWS announcements
func ThisMonth() (Announcements, error) {
	currentTime := time.Now()
	news, err := Fetch(currentTime.Year(), int(currentTime.Month()))
	if err != nil {
		return news, err
	}
	return news, nil
}

// Today get today's AWS announcments
func Today() (Announcements, error) {
	var todaysAnnouncements Announcements
	news, err := ThisMonth()
	if err != nil {
		return news, err
	}
	for _, announcement := range news {
		postDate, _ := time.Parse("Jan 2, 2006", announcement.PostDate)
		if dateEqual(postDate, time.Now()) {
			todaysAnnouncements = append(todaysAnnouncements, announcement)
		}
	}
	return todaysAnnouncements, nil
}

// Yesterday get yesterday's AWS announcments
func Yesterday() (Announcements, error) {
	var yesterdaysAnnouncements Announcements
	news, err := ThisMonth()
	if err != nil {
		return news, err
	}
	for _, announcement := range news {
		postDate, _ := time.Parse("Jan 2, 2006", announcement.PostDate)
		if dateEqual(postDate, time.Now().AddDate(0, 0, -1)) {
			yesterdaysAnnouncements = append(yesterdaysAnnouncements, announcement)
		}
	}
	return yesterdaysAnnouncements, nil
}

// Print Print out ASCII table of your selection of AWS announcements
func (a Announcements) Print() {
	data := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Announcement", "Date"})
	table.SetRowLine(true)
	for _, v := range a {
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
