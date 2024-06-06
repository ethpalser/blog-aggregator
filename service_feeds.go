package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ethpalser/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type RSS struct {
	XMLName    xml.Name   `xml:"rss"`
	RSSChannel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Items       []RSSData `xml:"item"`
}

type RSSData struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Url         string   `xml:"link"`
	Description string   `xml:"description"`
	PublishedAt string   `xml:"pubDate"`
}

type FeedService interface {
	getFeedsToFetch(since time.Time, n int32) (*[]FeedView, error)
	markFeedFetched(id uuid.UUID) error
	fetchFeed(url string) (*[]RSSData, error)
	processFeed(feed *RSSData) error
}

// Service-layer struct that contains business logic and dependencies
type RSSFeedService struct {
	DB *database.Queries
}

func NewFeedService(db *database.Queries) *RSSFeedService {
	return &RSSFeedService{
		DB: db,
	}
}

func (rfs *RSSFeedService) getFeedsToFetch(since time.Time, n int32) (*[]FeedView, error) {
	ctx := context.Background()
	dat, dbErr := rfs.DB.GetNextToFetchFeeds(ctx, database.GetNextToFetchFeedsParams{
		LastFetchedAt: sql.NullTime{Time: since, Valid: true},
		Limit:         n,
	})
	if dbErr != nil {
		return nil, errors.New("database error")
	}

	res := []FeedView{}
	for _, feed := range dat {
		res = append(res, DBFeedToView(&feed))
	}
	return &res, nil
}

func (rfs *RSSFeedService) markFeedFetched(id uuid.UUID) error {
	ctx := context.Background()
	_, dbErr := rfs.DB.UpdateFeedFetchedAt(ctx, uuid.NullUUID{UUID: id, Valid: true})
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func (rfs *RSSFeedService) fetchFeed(url string) (*[]RSSData, error) {
	// Perform https request at url
	resp, httpsErr := http.Get(url)
	if httpsErr != nil {
		return nil, fmt.Errorf("error fetching data: %s", httpsErr)
	}
	// Parse XML
	dat := RSS{}
	decoder := xml.NewDecoder(resp.Body)
	xmlErr := decoder.Decode(&dat)
	if xmlErr != nil {
		return nil, fmt.Errorf("error decoding xml: %s", xmlErr)
	}
	return &dat.RSSChannel.Items, nil
}

func (rfs *RSSFeedService) processFeed(feed *RSSData) error {
	println(feed.Title)
	return nil
}

func workerFetchFeeds(fs FeedService, fetchQuantity int) error {
	// Get feeds to fetch
	log.Printf("Fetching %v oldest feeds\n", fetchQuantity)
	toFetch, dbErr := fs.getFeedsToFetch(time.Now().Add(time.Hour*-1), int32(fetchQuantity))
	if dbErr != nil {
		log.Printf("Error fetching feching from database: %s\n", dbErr.Error())
		return dbErr
	}

	// Get data from feeds
	log.Printf("Fetching each feed asyncronously\n")
	wg := sync.WaitGroup{}
	for _, feed := range *toFetch {
		wg.Add(1)
		go func(f *FeedView) {
			defer wg.Done()

			log.Printf("Fetching from url: %s\n", f.Url)
			dat, rssErr := fs.fetchFeed(f.Url)
			if rssErr != nil {
				log.Printf("Error fetching from RSS feed: %s\n", rssErr.Error())
				return
			}

			log.Printf("Marking feeds as fetched from url: %s\n", f.Url)
			dbErr := fs.markFeedFetched(f.Id)
			if dbErr != nil {
				log.Printf("Error updating feed in database with id: %v", f.Id)
				return
			}

			log.Printf("Processing data returned from url: %s\n", f.Url)
			for _, post := range *dat {
				fs.processFeed(&post)
			}
		}(&feed)
	}

	// Wait for all HTTP fetches to complete and be processed
	wg.Wait()
	log.Println("Finished processing all feeds")
	return nil
}
