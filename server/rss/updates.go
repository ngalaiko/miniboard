package rss

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *Service) listenToUpdates(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
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
	return nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "rss",
	})
}
