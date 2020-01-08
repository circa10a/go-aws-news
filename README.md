# go-aws-news

Fetch what's new from AWS

## WIP

## Usage

```go
package main

import (
	"fmt"
	"time"

	awsNews "github.com/circa10a/go-aws-news"
)

func main() {
	currentTime := time.Now()
	news := awsNews.Fetch(currentTime.Year(), int(currentTime.Month()))
	for _, v := range news {
		fmt.Printf("Title: %v\n", v.Title)
		fmt.Printf("Link: %v\n", v.Link)
		fmt.Printf("Date: %v\n", v.PostDate)
		fmt.Println()
	}
}
```
