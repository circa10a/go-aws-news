package news

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fakePage describes one page of search results to serve, without callers
// having to care how the API shapes its JSON.
type fakePage struct {
	headlines []string
	totalHits int
}

// body renders the page the way the AWS search API would.
func (p fakePage) body() ([]byte, error) {
	items := make([]map[string]any, 0, len(p.headlines))
	for _, headline := range p.headlines {
		items = append(items, map[string]any{
			"item": map[string]any{
				"additionalFields": map[string]any{
					"headline":     headline,
					"headlineUrl":  "/about-aws/whats-new/2024/01/x",
					"postDateTime": "2024-01-15T00:00:00Z",
				},
			},
		})
	}

	return json.Marshal(map[string]any{
		"items": items,
		"metadata": map[string]any{
			"count":     len(p.headlines),
			"totalHits": p.totalHits,
		},
	})
}

// fakeTransport serves canned pages and records every request it receives, so
// tests can prove which HTTP client was used without touching the network.
type fakeTransport struct {
	onCall   func()
	pages    []fakePage
	requests []*http.Request
}

func (f *fakeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Real transports fail a request whose context is already done.
	if err := r.Context().Err(); err != nil {
		return nil, err
	}

	f.requests = append(f.requests, r)
	if f.onCall != nil {
		f.onCall()
	}

	page := fakePage{}
	if len(f.pages) > 0 {
		page = f.pages[0]
		f.pages = f.pages[1:]
	}

	body, err := page.body()
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(string(body))),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

func newFakeClient(pages ...fakePage) (*Client, *fakeTransport) {
	transport := &fakeTransport{pages: pages}
	return NewClient(WithHTTPClient(&http.Client{Transport: transport})), transport
}

func TestWithHTTPClientIsUsedForRequests(t *testing.T) {
	t.Parallel()
	client, transport := newFakeClient(fakePage{headlines: []string{"Announcement one"}, totalHits: 1})

	news, err := client.FetchYear(context.Background(), 2024)

	assert.NoError(t, err)
	assert.Len(t, news, 1)
	assert.Equal(t, "Announcement one", news[0].Title)
	assert.Len(t, transport.requests, 1)
	assert.Equal(t, "aws.amazon.com", transport.requests[0].URL.Host)
	assert.Equal(t, "whats-new-v2#year#2024", transport.requests[0].URL.Query().Get("tags.id"))
}

func TestWithBaseURLRedirectsRequests(t *testing.T) {
	t.Parallel()
	transport := &fakeTransport{pages: []fakePage{{headlines: []string{"one"}, totalHits: 1}}}
	client := NewClient(
		WithHTTPClient(&http.Client{Transport: transport}),
		WithBaseURL("https://example.test/custom/search"),
	)

	_, err := client.FetchYear(context.Background(), 2024)

	assert.NoError(t, err)
	assert.Len(t, transport.requests, 1)
	assert.Equal(t, "example.test", transport.requests[0].URL.Host)
	assert.True(t, strings.HasPrefix(transport.requests[0].URL.Path, "/custom/search"))
}

// WithHTTPClient replaces the underlying resty client, so it must not discard a
// base URL set by an earlier option.
func TestWithBaseURLSurvivesWithHTTPClient(t *testing.T) {
	t.Parallel()
	transport := &fakeTransport{pages: []fakePage{{headlines: []string{"one"}, totalHits: 1}}}
	client := NewClient(
		WithBaseURL("https://example.test/custom/search"),
		WithHTTPClient(&http.Client{Transport: transport}),
	)

	_, err := client.FetchYear(context.Background(), 2024)

	assert.NoError(t, err)
	assert.Len(t, transport.requests, 1)
	assert.Equal(t, "example.test", transport.requests[0].URL.Host)
}

func TestWithPostBaseURLPrefixesAnnouncementLinks(t *testing.T) {
	t.Parallel()
	transport := &fakeTransport{pages: []fakePage{{headlines: []string{"one"}, totalHits: 1}}}
	client := NewClient(
		WithHTTPClient(&http.Client{Transport: transport}),
		WithPostBaseURL("https://example.test"),
	)

	news, err := client.FetchYear(context.Background(), 2024)

	assert.NoError(t, err)
	assert.Len(t, news, 1)
	assert.Equal(t, "https://example.test/about-aws/whats-new/2024/01/x", news[0].Link)
}

func TestClientPaginatesUntilTotalHitsReached(t *testing.T) {
	t.Parallel()
	client, transport := newFakeClient(
		fakePage{headlines: []string{"one", "two"}, totalHits: 3},
		fakePage{headlines: []string{"three"}, totalHits: 3},
	)

	news, err := client.FetchYear(context.Background(), 2024)

	assert.NoError(t, err)
	assert.Len(t, news, 3)
	assert.Len(t, transport.requests, 2)
	assert.Equal(t, "0", transport.requests[0].URL.Query().Get("page"))
	assert.Equal(t, "1", transport.requests[1].URL.Query().Get("page"))
}

func TestClientReturnsErrorForCancelledContext(t *testing.T) {
	t.Parallel()
	client, _ := newFakeClient()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.FetchYear(ctx, 2024)

	assert.ErrorIs(t, err, context.Canceled)
}

// Cancellation has to take effect between pages, not just before the first
// request, otherwise a multi-page year ignores the caller's deadline.
func TestClientStopsPaginatingWhenContextCancelled(t *testing.T) {
	t.Parallel()
	client, transport := newFakeClient(
		fakePage{headlines: []string{"one", "two"}, totalHits: 9},
		fakePage{headlines: []string{"three"}, totalHits: 9},
	)
	ctx, cancel := context.WithCancel(context.Background())
	transport.onCall = cancel

	_, err := client.FetchYear(ctx, 2024)

	assert.ErrorIs(t, err, context.Canceled)
	assert.Len(t, transport.requests, 1)
}

func TestClientTodayUsesProvidedContext(t *testing.T) {
	t.Parallel()
	client, _ := newFakeClient()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Today(ctx)

	assert.ErrorIs(t, err, context.Canceled)
}
