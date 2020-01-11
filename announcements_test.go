package awsnews

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Integration test
func TestFetch(t *testing.T) {
	news, err := Fetch(2019, 12)
	assert.NoError(t, err)
	assert.Equal(t, len(news), 100)
}

func TestParseDate(t *testing.T) {
	assert.Equal(t, "Jan 1, 2020", parseDate("Posted on: Jan 1, 2020"))
}

// Integration Test
func TestThisMonth(t *testing.T) {
	today := time.Now()
	news, err := ThisMonth()
	assert.NoError(t, err)
	// Ensure each announcement returned matches current month
	for _, n := range news {
		assert.Equal(t, n.PostDate[:3], today.Month().String()[:3])
	}
}

func TestToday(t *testing.T) {
	today := time.Now()
	news, err := Today()
	assert.NoError(t, err)
	// Ensure each announcement returned matches current day
	for _, n := range news {
		postDate, _ := time.Parse("Jan 2, 2006", n.PostDate)
		assert.Equal(t, postDate.Day(), today.Day())
	}
}

func TestYesterday(t *testing.T) {
	news, err := Yesterday()
	assert.NoError(t, err)
	// Ensure each announcement returned matches yesterday
	for _, n := range news {
		postDate, _ := time.Parse("Jan 2, 2006", n.PostDate)
		assert.Equal(t, postDate.Day(), time.Now().AddDate(0, 0, -1).Day())
	}
}
