package feeds

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/ngalaiko/miniboard/server/genproto/feeds/v1"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
				log().Errorf("failed to update feeds: %s", err)
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

	wg, ctx := errgroup.WithContext(ctx)
	for _, feed := range feeds {
		feed := feed
		wg.Go(func() error {
			if err := s.updateFeed(ctx, feed); err != nil {
				return fmt.Errorf("failed to update feed: %w", err)
			}
			return nil
		})
	}

	return wg.Wait()
}

func (s *Service) updateFeed(ctx context.Context, feed *feeds.Feed) error {
	defer func() {
		if r := recover(); r != nil {
			log().Errorf("%s: %s", r, debug.Stack())
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

	if lastFetched.Add(updateInterval).After(time.Now()) {
		log().Infof("no need to update %s (%s): lastFetched at %s", feed.Id, feed.Url, lastFetched)
		return nil
	}

	resp, err := s.fetcher.Fetch(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", feed.Url, err)
	}
	defer resp.Body.Close()

	log().Infof("updating %s", feed.Id)

	items, err := s.parse(resp.Body, feed)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", feed.Url, err)
	}

	for _, item := range items {
		if err := s.saveItem(ctx, item, feed); err != nil {
			log().Errorf("failed to save item %s: %s", item.Link, err)
			continue
		}
	}

	if err := s.storage.Update(ctx, feed, feed.UserId); err != nil {
		return fmt.Errorf("failed to update feed: %w", err)
	}

	return nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "feeds",
	})
}
