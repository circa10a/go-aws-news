package news

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

// Client fetches AWS announcements. Use NewClient to create one.
type Client struct {
	resty       *resty.Client
	baseURL     string
	postBaseURL string
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sends requests through the given http.Client, which lets
// callers own transport concerns such as timeouts, proxies and instrumentation.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.resty = resty.NewWithClient(hc)
	}
}

// WithBaseURL queries the given search API instead of AWS's, which lets callers
// point at a proxy or a stub server.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithPostBaseURL builds announcement links against the given host. The search
// API returns only paths, so this is what makes a link complete.
func WithPostBaseURL(url string) Option {
	return func(c *Client) {
		c.postBaseURL = url
	}
}

// NewClient returns a Client configured with the given options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		resty:       resty.New(),
		baseURL:     awsWhatsNewBaseURL,
		postBaseURL: awsWhatsNewPostBaseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	// Applied last so an option swapping the underlying client cannot drop it.
	c.resty.SetBaseURL(c.baseURL)

	return c
}

// defaultClient backs the package level functions.
var defaultClient = NewClient()
