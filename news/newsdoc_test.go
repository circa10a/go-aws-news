package news

import (
	"log"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

// Reusable html documents
var monthTestDoc *goquery.Document
var yearTestDoc *goquery.Document

func init() {
	// create from a file
	// Monthly Testing
	monthFile, err := os.Open("./resources/aws-news-2019-12.html")
	if err != nil {
		log.Fatal(err)
	}
	defer monthFile.Close()
	monthDoc, err := goquery.NewDocumentFromReader(monthFile)
	if err != nil {
		//nolint
		log.Fatal(err)
	}
	monthTestDoc = monthDoc

	// create from a file
	// yearly Testing
	yearFile, err := os.Open("./resources/aws-news-2020.html")
	if err != nil {
		log.Fatal(err)
	}
	defer yearFile.Close()
	yearDoc, err := goquery.NewDocumentFromReader(yearFile)
	if err != nil {
		log.Fatal(err)
	}
	yearTestDoc = yearDoc
}

// Integration test
func TestGetNewsDocMonth(t *testing.T) {
	_, err := getNewsDocMonth(2019, 12)
	assert.NoError(t, err)
}

// Integration
func TestGetNewsDocYear(t *testing.T) {
	_, err := getNewsDocYear(2020)
	assert.NoError(t, err)
}

func TestGetSelectionItems(t *testing.T) {
	news := newsDoc{monthTestDoc}
	assert.Equal(t, news.getSelectionItems().Length(), 100)
}

func TestGetSelectionTitle(t *testing.T) {
	news := newsDoc{monthTestDoc}
	//nolint
	assert.Equal(t, news.getSelectionTitle(news.getSelectionItems().First()), "Amazon EKS Announces Beta Release of Amazon FSx for Lustre CSI Driver")
}

func TestGetSelectionLink(t *testing.T) {
	news := newsDoc{monthTestDoc}
	//nolint
	assert.Equal(t, news.getSelectionItemLink(news.getSelectionItems().First()), "https://aws.amazon.com/about-aws/whats-new/2019/12/amazon-eks-announces-beta-release-amazon-fsx-lustre-csi-driver/")
}

func TestGetSelectionDate(t *testing.T) {
	news := newsDoc{monthTestDoc}
	//nolint
	assert.Equal(t, news.getSelectionItemDate(news.getSelectionItems().First()), "\n              Posted On: Dec 23, 2019 \n            ")
}
