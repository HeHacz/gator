package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/hehacz/gator/internal/config"
)

func fetchFeed(ctx context.Context, feedURL string) (*config.RSSFeed, error) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	var rssFeed config.RSSFeed
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("coudnt create the request error: %v", err)
	}
	req.Header.Set("User-Agent", "Gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error coudnt read RSS feed: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error, coudnt read RSS feed body: %v", err)
	}
	xml.Unmarshal(body, &rssFeed)
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}
	return &rssFeed, nil
}
