package awsnews

import (
	"log"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

// Reusable html document
var testDoc *goquery.Document

func init() {
	// create from a file
	f, err := os.Open("./resources/aws-news-2019-12.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	testDoc = doc
}

// Integration test
func TestGetNewsDoc(t *testing.T) {
	_, err := getNewsDoc(2019, 12)
	assert.NoError(t, err)
}

func TestGetSelectionItems(t *testing.T) {
	news := newsDoc{testDoc}
	assert.Equal(t, news.getSelectionItems().Length(), 100)
}

func TestGetSelectionTitle(t *testing.T) {
	news := newsDoc{testDoc}
	assert.Equal(t, news.getSelectionTitle(news.getSelectionItems().First()), "Amazon EKS Announces Beta Release of Amazon FSx for Lustre CSI Driver")
}

func TestGetSelectionLink(t *testing.T) {
	news := newsDoc{testDoc}
	assert.Equal(t, news.getSelectionItemLink(news.getSelectionItems().First()), "https://aws.amazon.com/about-aws/whats-new/2019/12/amazon-eks-announces-beta-release-amazon-fsx-lustre-csi-driver/")
}

func TestGetSelectionDate(t *testing.T) {
	news := newsDoc{testDoc}
	assert.Equal(t, news.getSelectionItemDate(news.getSelectionItems().First()), "\n              Posted On: Dec 23, 2019 \n            ")
}
