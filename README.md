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
	// Most recent news
	currentNews := awsNews.ThisMonth()
	currentNews.Print()

	// Today's news
	todaysNews := awsNews.Today()
	todaysNews.Print()
}
```
