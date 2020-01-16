package providers

import (
	"fmt"
	"io/ioutil"

	"go-aws-news/news"
)

// Config is the raw data read from the provider configuration file.
var Config = readConfig()

func readConfig() []byte {
	fmt.Println("readConfig")
	b, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	return b
}

// Provider is implemented in each webhook provider in providers/.
type Provider interface {
	Notify(news.Announcements)
	GetName() string
	Enabled() bool
}

// providers is a slice of registered providers.
var providers []Provider

// GetProviders returns a list of registered providers.
func GetProviders() []Provider {
	return providers
}

// RegisterProvider adds the provider to the list of registered providers if enabled in the config.
func RegisterProvider(p Provider) {
	if p.Enabled() {
		providers = append(providers, p)
	}
}
