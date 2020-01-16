package news

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateEqual(t *testing.T) {
	today := time.Now()
	fakePostDateString := fmt.Sprintf("%v %02d, %v", today.Month().String()[:3], today.Day(), today.Year())
	fakePostDateTime, _ := time.Parse("Jan 2, 2006", fakePostDateString)
	assert.True(t, dateEqual(fakePostDateTime, today))
}
