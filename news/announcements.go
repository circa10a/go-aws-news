package news

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

const (
	awsWhatsNewBaseURL     = "https://aws.amazon.com/api/dirs/items/search"
	awsWhatsNewPostBaseURL = "https://aws.amazon.com"
	// awsWhatsNewPageSize is the largest page the search API will serve; asking
	// for more returns no results at all.
	awsWhatsNewPageSize = 2000
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

// getItemsYear gets every announcement AWS has tagged with the given year.
// A year holds more announcements than one page can carry, and results are
// sorted newest first, so a single request drops the earliest months entirely.
func (c *Client) getItemsYear(ctx context.Context, year int) (*AWSNewsItemsResponse, error) {
	results := &AWSNewsItemsResponse{}

	for page := 0; ; page++ {
		// Checked here so a cancellation between pages stops the walk instead of
		// issuing a request that is already doomed.
		if err := ctx.Err(); err != nil {
			return results, err
		}

		pageResults := &AWSNewsItemsResponse{}

		resp, err := c.resty.R().
			SetContext(ctx).
			SetResult(pageResults).
			SetQueryParams(map[string]string{
				"size":             strconv.Itoa(awsWhatsNewPageSize),
				"page":             strconv.Itoa(page),
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
			return results, fmt.Errorf("received response code: %d", resp.StatusCode())
		}

		results.FieldTypes = pageResults.FieldTypes
		results.Metadata = pageResults.Metadata
		results.Items = append(results.Items, pageResults.Items...)

		if len(pageResults.Items) == 0 || len(results.Items) >= pageResults.Metadata.TotalHits {
			break
		}
	}

	results.Metadata.Count = len(results.Items)

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
func (c *Client) Fetch(ctx context.Context, year int, month int) (Announcements, error) {
	announcements := Announcements{}
	items, err := c.getItemsYear(ctx, year)
	if err != nil {
		return Announcements{}, err
	}

	for _, item := range items.Items {
		announcement := Announcement{}
		// AWS's year tag does not always agree with postDateTime, so the year has
		// to be checked alongside the month to keep neighbouring years out.
		postDateYear, postDateMonth, _ := item.Item.AdditionalFields.PostDateTime.Date()
		if postDateYear == year && postDateMonth == time.Month(month) {
			announcement.Link = c.postBaseURL + item.Item.AdditionalFields.HeadlineURL
			announcement.PostDate = item.Item.AdditionalFields.PostDateTime.Format(time.RFC3339)
			announcement.Title = item.Item.AdditionalFields.Headline
			announcements = append(announcements, announcement)
		}
	}

	return announcements, nil
}

// FetchYear gets all of the announcements for the specified year that was input.
func (c *Client) FetchYear(ctx context.Context, year int) (Announcements, error) {
	announcements := Announcements{}
	items, err := c.getItemsYear(ctx, year)
	if err != nil {
		return announcements, err
	}

	for _, item := range items.Items {
		announcement := Announcement{}
		announcement.Link = c.postBaseURL + item.Item.AdditionalFields.HeadlineURL
		announcement.PostDate = item.Item.AdditionalFields.PostDateTime.Format(time.RFC3339)
		announcement.Title = item.Item.AdditionalFields.Headline
		announcements = append(announcements, announcement)
	}

	return announcements, nil
}

// fetchDay gets the announcements posted on the given date. The query covers the
// month the date itself falls in, which matters whenever the date belongs to a
// different month or year than today.
func (c *Client) fetchDay(ctx context.Context, date time.Time) (Announcements, error) {
	dayAnnouncements := Announcements{}

	items, err := c.Fetch(ctx, date.Year(), int(date.Month()))
	if err != nil {
		return dayAnnouncements, err
	}

	for _, announcement := range items {
		announcementPostDate, err := time.Parse(time.RFC3339, announcement.PostDate)
		if err != nil {
			return dayAnnouncements, err
		}

		if dateEqual(announcementPostDate, date) {
			dayAnnouncements = append(dayAnnouncements, announcement)
		}
	}

	return dayAnnouncements, nil
}

// ThisMonth gets the current month's AWS announcements.
func (c *Client) ThisMonth(ctx context.Context) (Announcements, error) {
	currentTime := time.Now()
	return c.Fetch(ctx, currentTime.Year(), int(currentTime.Month()))
}

// Today gets today's AWS announcements.
func (c *Client) Today(ctx context.Context) (Announcements, error) {
	return c.fetchDay(ctx, time.Now())
}

// Yesterday gets yesterday's AWS announcments.
func (c *Client) Yesterday(ctx context.Context) (Announcements, error) {
	return c.fetchDay(ctx, time.Now().AddDate(0, 0, -1))
}

// Fetch gets all of the announcements for the specified year/month that was input.
func Fetch(year int, month int) (Announcements, error) {
	return defaultClient.Fetch(context.Background(), year, month)
}

// FetchYear gets all of the announcements for the specified year that was input.
func FetchYear(year int) (Announcements, error) {
	return defaultClient.FetchYear(context.Background(), year)
}

// ThisMonth gets the current month's AWS announcements.
func ThisMonth() (Announcements, error) {
	return defaultClient.ThisMonth(context.Background())
}

// Today gets today's AWS announcements.
func Today() (Announcements, error) {
	return defaultClient.Today(context.Background())
}

// Yesterday gets yesterday's AWS announcments.
func Yesterday() (Announcements, error) {
	return defaultClient.Yesterday(context.Background())
}

// Print Prints out an ASCII table of your selection of AWS announcements.
func (a Announcements) Print() {
	data := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Announcement", "Date"})
	for _, v := range a {
		s := []string{
			v.Title,
			v.PostDate,
		}
		data = append(data, s)
	}
	for _, v := range data {
		_ = table.Append(v)
	}
	_ = table.Render()
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
				// An announcement matching several terms is still one announcement.
				break
			}
		}
	}
	return filteredAnnouncements
}

// HTML Converts Announcements to an unordered html list.
func (a Announcements) HTML() string {
	var list strings.Builder
	list.WriteString("<ul>")
	for _, v := range a {
		fmt.Fprintf(&list, "<li><a href=\"%v\">%v</a></li>", html.EscapeString(v.Link), html.EscapeString(v.Title))
	}
	list.WriteString("</ul>")
	return list.String()
}
