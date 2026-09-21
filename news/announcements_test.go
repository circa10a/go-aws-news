package news

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Integration test
func TestFetch(t *testing.T) {
	t.Parallel()
	news, err := Fetch(2019, 12)
	assert.NoError(t, err)
	assert.Greater(t, len(news), 100)
}

// Integration test
func TestFetchYear(t *testing.T) {
	t.Parallel()
	news, err := FetchYear(2020)
	assert.NoError(t, err)
	assert.Greater(t, len(news), 0)
}

// Integration Test
func TestThisMonth(t *testing.T) {
	t.Parallel()
	today := time.Now()
	news, err := ThisMonth()
	assert.NoError(t, err)
	// Ensure each announcement returned matches current month
	for _, n := range news {
		postDate, err := time.Parse(time.RFC3339, n.PostDate)
		assert.NoError(t, err)
		assert.Equal(t, postDate.Month(), today.Month())
	}
}

// Integration Test
func TestToday(t *testing.T) {
	t.Parallel()
	today := time.Now()
	news, err := Today()
	assert.NoError(t, err)
	// Ensure each announcement returned matches current day
	for _, n := range news {
		postDate, _ := time.Parse(time.RFC3339, n.PostDate)
		fmt.Println(n.PostDate)
		assert.Equal(t, postDate.Day(), today.Day())
	}
}

// Integration Test
func TestYesterday(t *testing.T) {
	t.Parallel()
	news, err := Yesterday()
	assert.NoError(t, err)
	// Ensure each announcement returned matches yesterday
	for _, n := range news {
		postDate, err := time.Parse(time.RFC3339, n.PostDate)
		assert.NoError(t, err)
		assert.Equal(t, postDate.Day(), time.Now().AddDate(0, 0, -1).Day())
	}
}

func Test5(t *testing.T) {
	t.Parallel()
	news, err := FetchYear(2020)
	news = news.Last(5)
	assert.NoError(t, err)
	assert.Equal(t, len(news), 5)
}

func Test1000(t *testing.T) {
	t.Parallel()
	news, err := FetchYear(2020)
	news = news.Last(100)
	assert.NoError(t, err)
	assert.Equal(t, len(news), 100)
}

func TestJSON(t *testing.T) {
	t.Parallel()
	news, err := FetchYear(2020)
	assert.NoError(t, err)
	_, jsonErr := news.JSON()
	assert.NoError(t, jsonErr)
}

func TestHTML(t *testing.T) {
	t.Parallel()
	news, err := FetchYear(2020)
	assert.NoError(t, err)
	data := []byte(news.HTML())
	err = xml.Unmarshal(data, new(interface{}))
	assert.NoError(t, err)
}

func TestFilter(t *testing.T) {
	t.Parallel()
	news, _ := FetchYear(2020)
	filteredNews := news.Filter([]string{"EKS", "ECS"})
	// 2020 has 2294 announcements; the previous expectation of 69 counted only the
	// 2000 that fit in a single unpaginated response.
	assert.Equal(t, len(filteredNews), 78)
}

// Integration test
// AWS caps the search API's page size, so a single request cannot cover a full
// year. Because results are sorted newest-first, the oldest months were the
// ones silently dropped.
func TestFetchOldestMonthOfCompletedYear(t *testing.T) {
	t.Parallel()
	news, err := Fetch(2024, 1)
	assert.NoError(t, err)
	assert.Greater(t, len(news), 100)
}

// Integration test
func TestFetchYearReturnsEveryItem(t *testing.T) {
	t.Parallel()
	news, err := FetchYear(2024)
	assert.NoError(t, err)
	assert.Greater(t, len(news), 2000)
}

// Integration test
// AWS's year tag does not strictly agree with postDateTime, so filtering on the
// month alone lets neighbouring years leak through.
func TestFetchOnlyReturnsRequestedYear(t *testing.T) {
	t.Parallel()
	news, err := Fetch(2025, 5)
	assert.NoError(t, err)
	assert.NotEmpty(t, news)
	for _, n := range news {
		postDate, err := time.Parse(time.RFC3339, n.PostDate)
		assert.NoError(t, err)
		assert.Equal(t, 2025, postDate.Year())
	}
}

// Integration test
func TestFetchYearLinksAreWellFormed(t *testing.T) {
	t.Parallel()
	news, err := FetchYear(2024)
	assert.NoError(t, err)
	assert.NotEmpty(t, news)
	for _, n := range news {
		assert.NotContains(t, strings.TrimPrefix(n.Link, "https://"), "//")
	}
}

func TestHTMLEscapesTextAndPreservesURL(t *testing.T) {
	t.Parallel()
	a := Announcements{{
		Title:    `Amazon S3 "Express" & <you>`,
		Link:     "https://aws.amazon.com/about-aws/whats-new/2026/01/foo",
		PostDate: "2026-01-01T00:00:00Z",
	}}

	html := a.HTML()

	assert.Contains(t, html, `href="https://aws.amazon.com/about-aws/whats-new/2026/01/foo"`)
	assert.Contains(t, html, "Amazon S3 &#34;Express&#34; &amp; &lt;you&gt;")
	assert.NotContains(t, html, "%2F")
}

func TestFilterDoesNotDuplicateMultipleMatches(t *testing.T) {
	t.Parallel()
	a := Announcements{{Title: "Amazon S3 bucket logging now GA"}}
	assert.Len(t, a.Filter([]string{"s3", "bucket"}), 1)
}

// Integration test
// Today/Yesterday used to query time.Now()'s year and month, so on the first of
// a month Yesterday() looked in the wrong month and always came back empty.
func TestFetchDayUsesTargetDateMonth(t *testing.T) {
	t.Parallel()
	date := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)

	news, err := defaultClient.fetchDay(context.Background(), date)

	assert.NoError(t, err)
	assert.NotEmpty(t, news)
	for _, n := range news {
		postDate, err := time.Parse(time.RFC3339, n.PostDate)
		assert.NoError(t, err)
		assert.True(t, dateEqual(postDate, date))
	}
}
