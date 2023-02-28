package news

import (
	"encoding/xml"
	"fmt"
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
	assert.Equal(t, len(filteredNews), 69)
}
