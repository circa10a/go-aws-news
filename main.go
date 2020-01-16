package main

import (
	"go-aws-news/providers"
	"log"

	"go-aws-news/news"
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
