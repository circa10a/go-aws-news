package rocketchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/circa10a/go-aws-news/providers"

	"github.com/circa10a/go-aws-news/news"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type config struct {
	Providers struct {
		Provider Provider `yaml:"rocketchat"`
	} `yaml:"providers"`
}

// Provider is an implementation of the `go-aws-news/providers` Provider interface.
type Provider struct {
	IsEnabled  bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhookURL"`
	IconURL    string `yaml:"iconURL"`
}

// Content is sent as json to the Rocket.Chat webhook endpoint.
type Content struct {
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
	Text     string `json:"text"`
}

// init initializes the provider from the provided config.
func init() {
	var c config
	if err := yaml.Unmarshal(providers.Config, &c); err != nil {
		log.Fatal(err)
	}

	providers.RegisterProvider(&c.Providers.Provider)
}

// Enabled returns true if the provider is enabled in the config.
func (p *Provider) Enabled() bool {
	return p.IsEnabled
}

// GetName returns the Provider's name.
func (*Provider) GetName() string {
	return "rocketchat"
}

// Notify is the function executed to POST to a provider's webhook url.
func (p *Provider) Notify(news news.Announcements) {
	var b strings.Builder
	for _, v := range news {
		b.WriteString(fmt.Sprintf("[%s](%s) - %s\n", v.Title, v.Link, v.PostDate))
	}

	content := &Content{
		Username: "AWS News",
		IconURL:  p.IconURL,
		Text:     b.String(),
	}

	json, err := json.Marshal(content)
	if err != nil {
		log.Error(fmt.Sprintf("[%s] %v", p.GetName(), err))
	}

	log.Info(fmt.Sprintf("[%v] Firing notification", p.GetName()))
	res, err := http.Post(p.WebhookURL, "application/json", bytes.NewBuffer(json))
	if err != nil {
		log.Error(fmt.Sprintf("[%s] %v", p.GetName(), err))
	}
	defer res.Body.Close()
}
