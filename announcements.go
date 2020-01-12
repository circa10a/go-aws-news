package awsnews

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
)

// Announcements Represents a slice containing all of the AWS announcements for a given time period.
type Announcements []Announcement

// Announcement Represents a single AWS product/feature announcement.
type Announcement struct {
	Title    string
	Link     string
	PostDate string
}

// newDoc Extends the goquery.Document to create more receiver functions on the *goquery.Document type.
type newsDoc struct {
	*goquery.Document
}

// Fetch Fetches all of the announcements for the specified year/month that was input.
func Fetch(year int, month int) (Announcements, error) {
	// Create []Announcement
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
		announcements = append(announcements, Announcement{Title: title, Link: link, PostDate: date})
	})
	return announcements, nil
}

// Extract date from amazon's format
// Input: Posted on: Jan 7, 2020
// Output: Jan 7, 2020
// parseDate Extracts a standarized date format from the AWS html document.
func parseDate(postDate string) string {
	r, _ := regexp.Compile("[A-Z][a-z]{2}\\s[0-9]{1,2},\\s[0-9]{4}")
	return r.FindStringSubmatch(postDate)[0]
}

// ThisMonth gets the current month's AWS announcements.
func ThisMonth() (Announcements, error) {
	currentTime := time.Now()
	news, err := Fetch(currentTime.Year(), int(currentTime.Month()))
	if err != nil {
		return news, err
	}
	return news, nil
}

// Today gets today's AWS announcements.
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

// Yesterday gets yesterday's AWS announcments.
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

// Print Prints out an ASCII table of your selection of AWS announcements.
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

// JSON Converts Announcements to JSON.
func (a Announcements) JSON() ([]byte, error) {
	json, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return json, nil
}
