package images

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
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

// SaveURL downloads and saves an image.
func (s *Service) SaveURL(ctx context.Context, articleName *resource.Name, url string) (*resource.Name, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download image")
	}
	return s.Save(ctx, articleName, resp.Body)
}

// Save saves an image.
func (s *Service) Save(ctx context.Context, articleName *resource.Name, reader io.Reader) (*resource.Name, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read image")
	}

	name := articleName.Child("images", ksuid.New().String())
	if err := s.storage.Store(ctx, name, data); err != nil {
		return nil, errors.Wrap(err, "failed to save image")
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
		w.Write(data)
	})
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
