# go-aws-news

Fetch what's new from AWS

## WIP

## Usage

```go
package main

import (
	awsNews "github.com/circa10a/go-aws-news"
)

func main() {
	// Today's news
	awsNews.Today().Print()
	// Yesterday's news
	awsNews.Yesterday.Print()
	// Current news
	awsNews.ThisMonth().Print()
	// Custom timeframe(June 2019)
	awsNews.Fetch(2020, 01).Print()
	// Console output
	// +--------------------------------+-------------+
	// |             TITLE              |    DATE     |
	// +--------------------------------+-------------+
	// | Amazon Translate introduces    | Jan 8, 2020 |
	// | Batch Translation              |             |
	// +--------------------------------+-------------+
	// | AWS Systems Manager Quick      | Jan 8, 2020 |
	// | Setup now supports targeting   |             |
	// | all instances                  |             |
	// +--------------------------------+-------------+

	// Loop slice of stucts of announcements
	// For your own data manipulation
	news := awsNews.Fetch(time.Now().Year(), int(time.Now().Month()))
	for _, v := range news {
		fmt.Printf("Title: %v\n", v.Title)
		fmt.Printf("Link: %v\n", v.Link)
		fmt.Printf("Date: %v\n", v.PostDate)
	}
}
```

## Table Output
