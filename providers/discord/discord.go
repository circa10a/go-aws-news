package discord

import (
	"fmt"
	"github.com/circa10a/go-aws-news/providers"
	"net/http"
	"net/url"
	"strings"

	"github.com/circa10a/go-aws-news/news"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type config struct {
	Providers struct {
		Provider Provider `yaml:"discord"`
	} `yaml:"providers"`
}

// Provider is an implementation of the `go-aws-news/providers` Provider interface.
type Provider struct {
	IsEnabled  bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhookURL"`
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
	return "discord"
}

// Notify is the function executed to POST to a provider's webhook url.
func (p *Provider) Notify(news news.Announcements) {

	var b strings.Builder
	for _, v := range news {
		b.WriteString(fmt.Sprintf("**Title:** *%v*\n**Link:** %v\n**Date:** %v\n", v.Title, v.Link, v.PostDate))
	}

	log.Info(fmt.Sprintf("[%v] Firing notification", p.GetName()))
	res, err := http.PostForm(p.WebhookURL, url.Values{
		"username": {"AWS News"},
		"content":  {b.String()},
	})
	if err != nil {
		log.Error(fmt.Sprintf("[%v] %v", p.GetName(), err))
	}
	defer res.Body.Close()
}
