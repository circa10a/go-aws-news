package providers

import (
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/circa10a/go-aws-news/news"
	log "github.com/sirupsen/logrus"
)

const defaultConfigName = "go-aws-news-config"

// Config is the raw data read from the provider configuration.
var Config []byte

func init() {
	var err error
	if os.Getenv("AWS_EXECUTION_ENV") == "AWS_Lambda_go1.x" {
		Config, err = lookupConfig()
	} else {
		Config, err = readConfigFile()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func lookupConfig() ([]byte, error) {
	configName := defaultConfigName
	if val, ok := os.LookupEnv("GO_AWS_NEWS_CONFIG_NAME"); ok {
		configName = val
	}

	client := ssm.New(session.Must(session.NewSession()))
	res, err := client.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(configName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return []byte(aws.StringValue(res.Parameter.Value)), nil
}

func readConfigFile() ([]byte, error) {
	b, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	return b, nil
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
