package yammer

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
		Provider Provider `yaml:"yammer"`
	} `yaml:"providers"`
}

// Provider is an implementation of the `go-aws-news/providers` Provider interface.
type Provider struct {
	IsEnabled bool   `yaml:"enabled"`
	APIURL    string `yaml:"APIURL"`
	GroupID   int    `yaml:"groupID"`
	Token     string `yaml:"token"`
}

// Payload is the respresentation of the json body being sent to Yammer via POST
type Payload struct {
	Body        string `json:"body"`
	GroupID     int    `json:"group_id"`
	IsRichText  bool   `json:"is_rich_text"`
	MessageType string `json:"message_type"`
	Title       string `json:"title"`
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
	return "yammer"
}

// UnescapedJSON creates its own encoder to ensure html characters are not converted to unicode
func (p Payload) UnescapedJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(p)
	return buffer.Bytes(), err
}

// Notify is the function executed to POST to a provider's webhook url.
func (p *Provider) Notify(news news.Announcements) {
	// Create payload
	payload := &Payload{
		Body:        news.HTML(),
		GroupID:     p.GroupID,
		IsRichText:  true,
		MessageType: "announcement",
		Title:       "AWS News",
	}
	json, _ := payload.UnescapedJSON()

	// Create request
	req, err := http.NewRequest("POST", p.APIURL, bytes.NewBuffer(json))
	if err != nil {
		log.Error(fmt.Sprintf("[%v] %v", p.GetName(), err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", p.Token))

	// Send request
	client := &http.Client{}
	log.Info(fmt.Sprintf("[%v] Firing notification", p.GetName()))
	resp, err := client.Do(req)
	if err != nil {
		log.Error(fmt.Sprintf("[%v] %v", p.GetName(), err))
	}
	// Cleanup
	defer resp.Body.Close()
}
