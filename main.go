package main

import (
	"github.com/circa10a/go-aws-news/providers"
	"log"

	"github.com/circa10a/go-aws-news/news"
)

func main() {

	news, err := news.Yesterday()
	if err != nil {
		log.Fatal(err)
	}

	providers := providers.GetProviders()

	for _, p := range providers {
		p.Notify(news)
	}
}
