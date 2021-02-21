package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/circa10a/go-aws-news/providers"
	log "github.com/sirupsen/logrus"

	"github.com/circa10a/go-aws-news/news"
)

func main() {
	if os.Getenv("AWS_EXECUTION_ENV") == "AWS_Lambda_go1.x" {
		lambda.Start(mainExec)
	} else {
		err := mainExec()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func mainExec() error {
	news, err := news.Yesterday()
	if err != nil {
		return err
	}
	if len(news) == 0 {
		log.Info("No news fetched. Skipping notifications.")
		return nil
	}

	providers := providers.GetProviders()

	if len(providers) == 0 {
		return errors.New("no providers enabled. see configuration")
	}
	for _, p := range providers {
		log.Info(fmt.Sprintf("[%v] Provider registered", p.GetName()))
		p.Notify(news)
	}

	return nil
}
