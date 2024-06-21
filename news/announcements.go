package news

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
)

const (
	awsWhatsNewBaseURL     = "https://aws.amazon.com/api/dirs/items/search"
	awsWhatsNewPostBaseURL = "https://aws.amazon.com"
)

type AWSNewsItemsResponse struct {
	FieldTypes struct {
		RelatedBlog  string `json:"relatedBlog"`
		PostBody     string `json:"postBody"`
		ModifiedDate string `json:"modifiedDate"`
		HeadlineURL  string `json:"headlineUrl"`
		PostDateTime string `json:"postDateTime"`
		PostSummary  string `json:"postSummary"`
		Headline     string `json:"headline"`
		ContentType  string `json:"contentType"`
	} `json:"fieldTypes"`
	Items []struct {
		Item struct {
			AdditionalFields struct {
				PostBody     string    `json:"postBody"`
				ModifiedDate time.Time `json:"modifiedDate"`
				HeadlineURL  string    `json:"headlineUrl"`
				PostDateTime time.Time `json:"postDateTime"`
				PostSummary  string    `json:"postSummary"`
				ContentType  string    `json:"contentType"`
				Headline     string    `json:"headline"`
			} `json:"additionalFields"`
			ID             string `json:"id"`
			Locale         string `json:"locale"`
			DirectoryID    string `json:"directoryId"`
			Name           string `json:"name"`
			CreatedBy      string `json:"createdBy"`
			LastUpdatedBy  string `json:"lastUpdatedBy"`
			DateCreated    string `json:"dateCreated"`
			DateUpdated    string `json:"dateUpdated"`
			Author         string `json:"author"`
			NumImpressions int    `json:"numImpressions"`
		} `json:"item"`
		Tags []struct {
			ID             string `json:"id"`
			Locale         string `json:"locale"`
			TagNamespaceID string `json:"tagNamespaceId"`
			Name           string `json:"name"`
			Description    string `json:"description"`
			CreatedBy      string `json:"createdBy"`
			LastUpdatedBy  string `json:"lastUpdatedBy"`
			DateCreated    string `json:"dateCreated"`
			DateUpdated    string `json:"dateUpdated"`
		} `json:"tags"`
	} `json:"items"`
	Metadata struct {
		Count     int `json:"count"`
		TotalHits int `json:"totalHits"`
	} `json:"metadata"`
}

func getItemsYear(year int) (*AWSNewsItemsResponse, error) {
	results := &AWSNewsItemsResponse{}
	client := resty.New()

	resp, err := client.SetBaseURL(awsWhatsNewBaseURL).R().
		SetResult(results).
		SetQueryParams(map[string]string{
			"size":             "2000", // 2000 seems to be the max or no results return
			"item.directoryId": "whats-new-v2",
			"sort_by":          "item.additionalFields.postDateTime",
			"sort_order":       "desc",
			"item.locale":      "en_US",
			"tags.id":          fmt.Sprintf("whats-new-v2#year#%d", year),
		}).
		SetHeader("Accept", "application/json").
		Get("/")

	if err != nil {
		return results, err
	}

	if resp.StatusCode() > 399 {
		return results, fmt.Errorf("Received response code: %d", resp.StatusCode())
	}

	return results, nil
}

// Announcements Represents a slice containing all of the AWS announcements for a given time period.
type Announcements []Announcement

// Announcement Represents a single AWS product/feature announcement.
type Announcement struct {
	Title    string
	Link     string
	PostDate string
}

// Fetch gets all of the announcements for the specified year/month that was input.
func Fetch(year int, month int) (Announcements, error) {
	announcements := Announcements{}
	items, err := getItemsYear(year)
	if err != nil {
		return Announcements{}, err
	}

	for _, item := range items.Items {
		announcement := Announcement{}
		_, postDateMonth, _ := item.Item.AdditionalFields.PostDateTime.Date()
		if postDateMonth == time.Now().Month() {
			announcement.Link = awsWhatsNewPostBaseURL + item.Item.AdditionalFields.HeadlineURL
			announcement.PostDate = item.Item.AdditionalFields.PostDateTime.Format(time.RFC3339)
			announcement.Title = item.Item.AdditionalFields.Headline
			announcements = append(announcements, announcement)
		}
	}

	return announcements, nil
}

// FetchYear gets all of the announcements for the specified year that was input.
func FetchYear(year int) (Announcements, error) {
	announcements := Announcements{}
	items, err := getItemsYear(year)
	if err != nil {
		return announcements, err
	}

	for _, item := range items.Items {
		announcement := Announcement{}
		announcement.Link = fmt.Sprintf("%s/%s", awsWhatsNewPostBaseURL, item.Item.AdditionalFields.HeadlineURL)
		announcement.PostDate = item.Item.DateCreated
		announcement.Title = item.Item.AdditionalFields.Headline
		announcements = append(announcements, announcement)
	}

	return announcements, nil
}

// ThisMonth gets the current month's AWS announcements.
func ThisMonth() (Announcements, error) {
	currentTime := time.Now()
	items, err := Fetch(currentTime.Year(), int(currentTime.Month()))
	if err != nil {
		return items, err
	}
	return items, nil
}

// Today gets today's AWS announcements.
func Today() (Announcements, error) {
	todaysAnnouncements := Announcements{}
	currentTime := time.Now()

	items, err := Fetch(currentTime.Year(), int(currentTime.Month()))
	if err != nil {
		return todaysAnnouncements, err
	}

	for _, announcement := range items {
		announcementPostDate, err := time.Parse(time.RFC3339, announcement.PostDate)
		if err != nil {
			return todaysAnnouncements, err
		}

		if dateEqual(announcementPostDate, currentTime) {
			todaysAnnouncements = append(todaysAnnouncements, announcement)
		}
	}

	return todaysAnnouncements, nil
}

// Yesterday gets yesterday's AWS announcments.
func Yesterday() (Announcements, error) {
	yesterdaysAnnouncements := Announcements{}
	currentTime := time.Now()
	yesterday := currentTime.AddDate(0, 0, -1)

	items, err := Fetch(currentTime.Year(), int(currentTime.Month()))
	if err != nil {
		return yesterdaysAnnouncements, err
	}

	for _, announcement := range items {
		announcementPostDate, err := time.Parse(time.RFC3339, announcement.PostDate)
		if err != nil {
			return yesterdaysAnnouncements, err
		}

		if dateEqual(announcementPostDate, yesterday) {
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
		html.WriteString(fmt.Sprintf("<li><a href='%v'>%v</a></li>", url.QueryEscape(v.Link), url.QueryEscape(v.Title)))
	}
	html.WriteString("</ul>")
	return html.String()
}
