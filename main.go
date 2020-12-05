package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circa10a/go-aws-news/providers"
	log "github.com/sirupsen/logrus"

	"github.com/circa10a/go-aws-news/news"
)

func main() {
	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		lambda.Start(mainExec)
	} else {
		mainExec()
	}
}

func mainExec() {
	news, err := news.Yesterday()
	if err != nil {
		log.Fatal(err)
	}
	if len(news) == 0 {
		log.Info("No news fetched. Skipping notifications.")
		log.Exit(0)
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
