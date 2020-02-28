package images

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"net/http"

	"image/jpeg"
	// supported image formats
	_ "image/png"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Service manages images.
type Service struct {
	storage storage.Storage
}

// New returns new service instance.
func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// Save saves an image.
func (s *Service) Save(ctx context.Context, articleName *resource.Name, reader io.Reader) (*resource.Name, error) {
	imageData, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	jpegImage := &bytes.Buffer{}
	if err := jpeg.Encode(jpegImage, imageData, nil); err != nil {
		return nil, fmt.Errorf("failed to encode jpeg: %w", err)
	}

	name := articleName.Child("images", ksuid.New().String())
	if err := s.storage.Store(ctx, name, jpegImage.Bytes()); err != nil {
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
