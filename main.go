package main

import (
	"fmt"

	"github.com/circa10a/go-aws-news/providers"
	log "github.com/sirupsen/logrus"

	"github.com/circa10a/go-aws-news/news"
)

func main() {

	news, err := news.Yesterday()
	if err != nil {
		log.Fatal(err)
	}

	providers := providers.GetProviders()

	if len(providers) == 0 {
		log.Fatal("No providers enabled. See config.yaml.")
	}

	for _, p := range providers {
		log.Info(fmt.Sprintf("[%v] Provider registered", p.GetName()))
		p.Notify(news)
	}
}
