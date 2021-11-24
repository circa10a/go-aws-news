package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/circa10a/go-aws-news/providers"

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
	WebhookURL string `yaml:"webhookURL"`
	AvatarURL  string `yaml:"avatarURL"`
	IsEnabled  bool   `yaml:"enabled"`
}

// Embed is a custom field for a Discord webhook payload.
type Embed struct {
	Description string `json:"description"`
	Color       string `json:"color"`
}

// Payload is the json sent to the Discord webhook endpoint.
type Payload struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []Embed `json:"embeds"`
}

const (
	awsOrange = "16750848"
	username  = "AWS News"
	maxEmbeds = 10
)

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
	var embeds []Embed

	for _, v := range news {
		description := fmt.Sprintf("[%s](%s)\n%s", v.Title, v.Link, v.PostDate)
		embeds = append(embeds, Embed{description, awsOrange})

		// Discord has a limit on num embeds per message
		if len(embeds) == maxEmbeds {
			err := p.post(embeds)
			if err != nil {
				log.Error(err)
			}

			embeds = nil
		}
	}

	if len(embeds) > 0 {
		err := p.post(embeds)
		if err != nil {
			log.Error(err)
		}
	}
}

func (p *Provider) post(embeds []Embed) error {
	payload := &Payload{
		Username:  username,
		AvatarURL: p.AvatarURL,
		Embeds:    embeds,
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[%s] %w", p.GetName(), err)
	}

	log.Info(fmt.Sprintf("[%v] Firing notification", p.GetName()))
	res, err := http.Post(p.WebhookURL, "application/json", bytes.NewBuffer(json))
	if err != nil {
		return fmt.Errorf("[%v] %w", p.GetName(), err)
	}
	defer res.Body.Close()

	return nil
}
