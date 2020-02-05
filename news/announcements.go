package news

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
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

func (n newsDoc) GetAnnouncements() (Announcements, error) {
	// Create []Announcement
	var announcements Announcements
	n.getSelectionItems().Each(func(i int, s *goquery.Selection) {
		title := n.getSelectionTitle(s)
		link := fmt.Sprintf("https:%v", n.getSelectionItemLink(s))
		date := parseDate(n.getSelectionItemDate(s))
		announcements = append(announcements, Announcement{Title: title, Link: link, PostDate: date})
	})
	return announcements, nil
}

// Fetch gets all of the announcements for the specified year/month that was input.
func Fetch(year int, month int) (Announcements, error) {
	doc, err := getNewsDocMonth(year, month)
	if err != nil {
		return Announcements{}, err
	}
	return newsDoc{doc}.GetAnnouncements()
}

// FetchYear gets all of the announcements for the specified year that was input.
func FetchYear(year int) (Announcements, error) {
	doc, err := getNewsDocYear(year)
	if err != nil {
		return Announcements{}, err
	}
	return newsDoc{doc}.GetAnnouncements()
}

// ThisMonth gets the current month's AWS announcements.
func ThisMonth() (Announcements, error) {
	var thisMonthsAnnouncements Announcements
	currentTime := time.Now()
	news, err := FetchYear(currentTime.Year())
	for _, announcement := range news {
		if announcement.PostDate[:3] == currentTime.Month().String()[:3] {
			thisMonthsAnnouncements = append(thisMonthsAnnouncements, announcement)
		}
	}
	if err != nil {
		return news, err
	}
	return thisMonthsAnnouncements, nil
}

// Today gets today's AWS announcements.
func Today() (Announcements, error) {
	var todaysAnnouncements Announcements
	news, err := FetchYear(time.Now().Year())
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
	news, err := FetchYear(time.Now().Year())
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

// Extract date from amazon's format
// Input: Posted on: Jan 7, 2020
// Output: Jan 7, 2020
// parseDate Extracts a standarized date format from the AWS html document.
func parseDate(postDate string) string {
	r, _ := regexp.Compile("[A-Z][a-z]{2}\\s[0-9]{1,2},\\s[0-9]{4}")
	// AWS sometimes doesn't have a post date
	if len(r.FindStringSubmatch(postDate)) > 0 {
		return r.FindStringSubmatch(postDate)[0]
	}
	return "No posted date"
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

// Last returns a set number of news items you specify
func (a Announcements) Last(n int) Announcements {
	if len(a) > n {
		return a[:n]
	}
	return a
}

// JSON Converts Announcements to JSON.
func (a Announcements) JSON() ([]byte, error) {
	json, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return json, nil
}

// Filter accepts a slice of products/terms to only return announcments you care about
func (a Announcements) Filter(p []string) Announcements {
	var filteredAnnouncements Announcements
	for _, v := range a {
		for _, product := range p {
			if strings.Contains(strings.ToLower(v.Title), strings.ToLower(product)) {
				filteredAnnouncements = append(filteredAnnouncements, v)
			}
		}
	}
	return filteredAnnouncements
}

// HTML Converts Announcements to an unordered html list.
func (a Announcements) HTML() string {
	var html strings.Builder
	html.WriteString("<ul>")
	for _, v := range a {
		html.WriteString(fmt.Sprintf("<li><a href='%v'>%v</a></li>", v.Link, v.Title))
	}
	html.WriteString("</ul>")
	return html.String()
}
