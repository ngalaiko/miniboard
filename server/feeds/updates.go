package feeds

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/storage/resource"
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
	raw, err := s.storage.LoadAll(ctx, resource.ParseName("users/*/feeds/*"))
	if err != nil {
		return fmt.Errorf("failed to load feeds: %w", err)
	}

	wg, ctx := errgroup.WithContext(ctx)
	for _, r := range raw {
		feed := &Feed{}
		if err := proto.Unmarshal(r.Data, feed); err != nil {
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

func (s *Service) updateFeed(ctx context.Context, feed *Feed) error {
	defer func() {
		if r := recover(); r != nil {
			log().Errorf("%s: %s", r, debug.Stack())
		}
	}()

	name := resource.ParseName(feed.Name)

	ctx = actor.NewContext(ctx, name.Parent())

	lastFetched := time.Time{}
	if feed.LastFetched != nil {
		var err error
		lastFetched, err = ptypes.Timestamp(feed.LastFetched)
		if err != nil {
			return fmt.Errorf("failed to unmarshal last_fetched: %w", err)
		}
	}

	if lastFetched.Add(updateInterval).After(time.Now()) {
		log().Infof("no need to update %s (%s): lastFetched at %s", feed.Name, feed.Url, lastFetched)
		return nil
	}

	resp, err := s.fetcher.Fetch(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", feed.Url, err)
	}
	defer resp.Body.Close()

	log().Infof("updating %s", feed.Name)

	if err := s.parse(ctx, resp.Body, feed); err != nil {
		return fmt.Errorf("failed to parse %s: %w", feed.Url, err)
	}

	return nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "feeds",
	})
}
