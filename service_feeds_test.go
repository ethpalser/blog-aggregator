package main

import (
	"testing"
	"time"

	"github.com/ethpalser/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type RSSFeedServiceMock struct {
	DB *database.Queries
}

func (fs *RSSFeedServiceMock) getFeedsToFetch(since time.Time, n int32) (*[]FeedView, error) {
	testCases := []FeedView{
		{
			Name: "Boot.dev Blog",
			Url:  "https://blog.boot.dev/index.xml",
		},
		{
			Name: "Wagslane.dev Blog",
			Url:  "https://wagslane.dev/index.xml",
		},
	}
	return &testCases, nil
}

func (fs *RSSFeedServiceMock) markFeedFetched(id uuid.UUID) error {
	// Do nothing
	return nil
}

func (fs *RSSFeedServiceMock) fetchFeed(url string) (*RSSData, error) {
	testResponse := map[string]RSSData{
		"https://blog.boot.dev/index.xml": {
			Title: "Behind the Scenes: Boots, the Greatest Companion!",
		},
		"https://wagslane.dev/index.xml": {
			Title: "Why I Started Boot.dev",
		},
	}
	res := testResponse[url]
	return &res, nil
}

func (fs *RSSFeedServiceMock) processFeed(feed *RSSData) error {
	println(feed.Title)
	return nil
}

// Test fetching and processing in parallel (smoke test)
func TestUnit_WorkerFetchFeeds(t *testing.T) {
	rfs := RSSFeedServiceMock{
		DB: nil,
	}

	err := workerFetchFeeds(&rfs, 10)
	if err != nil {
		assert.Fail(t, "error fetching and processing feeds: %s", err.Error())
	}
}

func TestIntegration_FetchFeedAndProcessFeed(t *testing.T) {
	rfs := RSSFeedService{
		DB: nil, // Not required for test
	}

	testCases := []FeedView{
		{
			Name: "Boot.dev Blog",
			Url:  "https://blog.boot.dev/index.xml",
		},
		{
			Name: "Wagslane.dev Blog",
			Url:  "https://wagslane.dev/index.xml",
		},
	}

	for _, feed := range testCases {
		dat, rssErr := rfs.fetchFeed(feed.Url)
		if rssErr != nil {
			assert.Fail(t, "Failed fetching data from test feeds: %s", rssErr.Error())
			return
		}
		err := rfs.processFeed(dat)
		if err != nil {
			assert.Fail(t, "Failed processing data of feed: %v, %s", dat, err.Error())
			return
		}
	}
}
