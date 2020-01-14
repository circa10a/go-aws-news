# go-aws-news

Fetch what's new from AWS

<p align="center"><img src="https://i.imgur.com/U7zlAGc.png"/></p>

[![Build Status](https://travis-ci.org/circa10a/go-aws-news.svg?branch=master)](https://travis-ci.org/circa10a/go-aws-news)
[![Go Doc](https://godoc.org/github.com/circa10a/go-aws-news?status.svg)](http://godoc.org/github.com/circa10a/go-aws-news)
[![Go Report Card](https://goreportcard.com/badge/github.com/circa10a/go-aws-news)](https://goreportcard.com/report/github.com/circa10a/go-aws-news)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/circa10a/go-aws-news?style=plastic)

[go-aws-news](#go-aws-news)
  * [Install](#install)
  * [Usage](#usage)
      - [Get Today's news](#get-today-s-news)
      - [Get Yesterday's news](#get-yesterday-s-news)
      - [Get all news for the month](#get-all-news-for-the-month)
      - [Gets from a previous month](#gets-from-a-previous-month)
      - [Print out announcements](#print-out-announcements)
      - [Loop over news data](#loop-over-news-data)
      - [Get news as JSON](#get-news-as-json)
      - [Get news about a specific product](#get-news-about-a-specific-product)
  * [Development](#development)
    + [Test](#test)

## Install

```shell
go get -u "github.com/circa10a/go-aws-news"
```

## Usage

Methods return a slice of structs which include the announcement title, a link, and the date it was posted as well an error. This allows you to manipulate the data in whichever way you please, or simply use `Print()` to print a nice ASCII table to the console.

#### Get Today's news

```go
package main

import "github.com/circa10a/go-aws-news"

func main() { 
	news, err := awsnews.Today()
	if err != nil {
		// Handle error
	}
	news.Print()
}
```

#### Get Yesterday's news

```go
news, _ := awsnews.Yesterday()
```

#### Get all news for the month

```go
news, _ := awsnews.ThisMonth()
```

#### Gets from a previous month

```go
// Custom timeframe(June 2019)
news, err := awsnews.Fetch(2019, 06)
```

#### Print out announcements

```go
news, _ := awsnews.ThisMonth()
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
```

#### Loop over news data

```go
// Loop slice of stucts of announcements
// For your own data manipulation
news, _ := awsnews.Fetch(time.Now().Year(), int(time.Now().Month()))
for _, v := range news {
	fmt.Printf("Title: %v\n", v.Title)
	fmt.Printf("Link: %v\n", v.Link)
	fmt.Printf("Date: %v\n", v.PostDate)
}
```

#### Get news as JSON

```go
news, _ := awsnews.ThisMonth()
json, jsonErr := news.JSON()
if jsonErr != nil {
	log.Fatal(err)
}
fmt.Println(string(json))
```

#### Get news about a specific product

```go
news, err := awsnews.Fetch(2019, 12)
if err != nil {
	fmt.Println(err)
} else {
	news.Filter([]string{"EKS", "ECS"}).Print()
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

