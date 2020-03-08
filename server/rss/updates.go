package rss

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"miniboard.app/api/actor"
	rss "miniboard.app/proto/users/rss/v1"
	"miniboard.app/storage/resource"
)

const (
	pollInterval   = time.Minute
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
			if err := s.updateFeeds(ctx); err != nil {
				log().Errorf("failed to update feeds: %s", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) updateFeeds(ctx context.Context) error {
	raw, err := s.storage.LoadAll(ctx, resource.ParseName("users/*/rss/*"))
	if err != nil {
		return fmt.Errorf("failed to load feeds: %w", err)
	}

	wg, ctx := errgroup.WithContext(ctx)
	for _, r := range raw {
		feed := &rss.Feed{}
		if err := proto.Unmarshal(r, feed); err != nil {
			return fmt.Errorf("failed to unmarshal feed: %w", err)
		}

		wg.Go(func() error {
			if err := s.updateFeed(ctx, feed); err != nil {
				return fmt.Errorf("failed to update feed: %w", err)
			}
			return nil
		})
	}

	return wg.Wait()
}

func (s *Service) updateFeed(ctx context.Context, feed *rss.Feed) error {
	defer func() {
		if r := recover(); r != nil {
			log().Panicf("%s: %s", r, debug.Stack())
		}
	}()

	name := resource.ParseName(feed.Name)

	ctx = actor.NewContext(ctx, name.Parent())

	lastFetched, err := ptypes.Timestamp(feed.LastFetched)
	if err != nil {
		return fmt.Errorf("failed to unmarshal last_fetched: %w", err)
	}

	if lastFetched.Add(updateInterval).After(time.Now()) {
		return nil
	}

	resp, err := http.Get(feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", feed.Url, err)
	}
	defer resp.Body.Close()

	log().Infof("updating %s", feed.Name)

	if err := s.parse(ctx, resp.Body); err != nil {
		return fmt.Errorf("failed to parse %s: %w", feed.Url, err)
	}

	feed.LastFetched = ptypes.TimestampNow()

	raw, err := proto.Marshal(feed)
	if err != nil {
		return fmt.Errorf("failed to marshal feed: %w", err)
	}

	if err := s.storage.Store(ctx, name, raw); err != nil {
		return fmt.Errorf("failed to save feed: %w", err)
	}

	return nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "rss",
	})
}
