package images

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spaolacci/murmur3"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Service manages images.
type Service struct {
	storage storage.Storage
}

// NewService returns new service instance.
func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// Save saves an image.
func (s *Service) Save(ctx context.Context, reader io.Reader) (*resource.Name, error) {
	imageData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	hash := murmur3.New128()
	_, _ = hash.Write(imageData)

	name := resource.NewName("images", fmt.Sprintf("%x", hash.Sum(nil)))

	if err := s.storage.Store(ctx, name, imageData); err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	return name, nil
}

// Handler implements http.HandlerFunc.
func (s *Service) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := resource.ParseName(r.RequestURI[1:])

		data, err := s.storage.Load(r.Context(), name)
		if err != nil {
			log("images").Errorf("failed to load image %s: %s", name, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", http.DetectContentType(data))
		_, _ = w.Write(data)
	})
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
