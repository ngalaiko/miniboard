package feeds

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/ngalaiko/miniboard/server/genproto/feeds/v1"
)

const (
	pollInterval   = 5 * time.Minute
	updateInterval = 5 * time.Minute
)

func (s *Service) listenToUpdates(ctx context.Context) {
	ticker := time.NewTicker(pollInterval)
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithCancel(context.Background())
			if err := s.updateFeeds(ctx); err != nil {
				s.logger.Error("failed to update feeds: %s", err)
			}
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) updateFeeds(ctx context.Context) error {
	feeds, err := s.storage.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load feeds: %w", err)
	}

	for _, feed := range feeds {
		if err := s.updateFeed(ctx, feed); err != nil {
			s.logger.Error("failed to update feed: %s", err)
		}
	}

	return nil
}

func (s *Service) updateFeed(ctx context.Context, feed *feeds.Feed) error {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("%s: %s", r, debug.Stack())
		}
	}()

	lastFetched := time.Time{}
	if feed.LastFetched != nil {
		var err error
		lastFetched, err = ptypes.Timestamp(feed.LastFetched)
		if err != nil {
			return fmt.Errorf("failed to unmarshal last_fetched: %w", err)
		}
	}

	if lastFetched.Add(updateInterval).After(s.nowFunc()) {
		s.logger.Info("no need to update %s (%s): lastFetched at %s", feed.Id, feed.Url, lastFetched)
		return nil
	}

	resp, err := s.fetcher.Fetch(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", feed.Url, err)
	}
	defer resp.Body.Close()

	s.logger.Info("updating %s", feed.Id)

	items, err := s.parse(resp.Body, feed)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", feed.Url, err)
	}

	for _, item := range items {
		if err := s.saveItem(ctx, item, feed); err != nil {
			s.logger.Error("failed to save item %s: %s", item.Link, err)
			continue
		}
	}

	if err := s.storage.Update(ctx, feed, feed.UserId); err != nil {
		return fmt.Errorf("failed to update feed: %w", err)
	}

	return nil
}
