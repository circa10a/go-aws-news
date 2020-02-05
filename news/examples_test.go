package news

import "fmt"

func ExampleFetch() {
	news, err := Fetch(2019, 12)
	if err != nil {
		// Handle error
	}
	news.Print()
}

func ExampleFetchYear() {
	news, err := FetchYear(2020)
	if err != nil {
		// Handle error
	}
	news.Print()
}

func ExampleThisMonth() {
	news, err := ThisMonth()
	if err != nil {
		// Handle error
	}
	for _, n := range news {
		fmt.Println(n.Title)
	}
}

func ExampleToday() {
	news, err := Today()
	if err != nil {
		// Handle error
	}
	fmt.Println(news)
}

func ExampleYesterday() {
	news, err := Yesterday()
	if err != nil {
		// Handle error
	}
	fmt.Println(news)
}

func (a Announcements) ExampleJSON() {
	news, err := ThisMonth()
	if err != nil {
		// Handle error
	}
	json, jsonErr := news.Filter([]string{"S3"}).JSON()
	if jsonErr != nil {
		// Handle error
	}
	fmt.Println(string(json))
}

func (a Announcements) ExampleHTML() {
	news, err := ThisMonth()
	if err != nil {
		// Handle error
	}
	fmt.Println(news.Filter([]string{"S3"}).HTML())
}

func (a Announcements) ExampleFilter() {
	news, err := ThisMonth()
	if err != nil {
		// Handle error
	}
	fmt.Println(news.Filter([]string{"S3"}))
}
