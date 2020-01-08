# go-aws-news

Fetch what's new from AWS

## WIP

## Usage

```go
package main

import (
	"time"

	aws "github.com/circa10a/go-aws-news"
)

func main() {
	currentTime := time.Now()
	news := aws.News{
		Year:  currentTime.Year(),
		Month: int(currentTime.Month()),
	}
	news.Print()
}
```
