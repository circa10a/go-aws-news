# go-aws-news

Fetch what's new from AWS

<p align="center"><img src="https://i.imgur.com/U7zlAGc.png"/></p>

[![Build Status](https://travis-ci.org/circa10a/go-aws-news.svg?branch=master)](https://travis-ci.org/circa10a/go-aws-news)
[![Go Doc](https://godoc.org/github.com/circa10a/go-aws-news?status.svg)](http://godoc.org/github.com/circa10a/go-aws-news)
[![Go Report Card](https://goreportcard.com/badge/github.com/circa10a/go-aws-news)](https://goreportcard.com/report/github.com/circa10a/go-aws-news)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/circa10a/go-aws-news?style=plastic)

## Usage

Methods return a slice of structs which include the announcement title, a link, and the date it was posted as well an error. This allows you to manipulate the data in whichever way you please, or simply use `Print()` to print a nice ASCII table to the console.

```go
package main

import (
	awsNews "github.com/circa10a/go-aws-news"
)

func main() {
	// Today's news
	news, _ := awsNews.Today()
	news.Print()
	// Yesterday's news
	news, _:= awsNews.Yesterday()
	news.Print()
	// Current news
	news, err := awsNews.ThisMonth()
	news.Print()
	// Custom timeframe(June 2019)
	news, err := awsNews.Fetch(2020, 01)
	news.Print()
	// Console output
	// +--------------------------------+--------------+
	// |          ANNOUNCEMENT          |     DATE     |
	// +--------------------------------+--------------+
	// | Amazon Cognito now supports    | Jan 10, 2020 |
	// | CloudWatch Usage Metrics       |              |
	// +--------------------------------+--------------+
	// | Introducing Workload Shares in | Jan 10, 2020 |
	// | AWS Well-Architected Tool      |              |
	// +--------------------------------+--------------+
	//
	// Loop slice of stucts of announcements
	// For your own data manipulation
	news, _ := awsNews.Fetch(time.Now().Year(), int(time.Now().Month()))
	for _, v := range news {
		fmt.Printf("Title: %v\n", v.Title)
		fmt.Printf("Link: %v\n", v.Link)
		fmt.Printf("Date: %v\n", v.PostDate)
	}
}
```

## Development

### Test

```shell
# Unit/Integration tests
make
# Get code coverage
make coverage
```

