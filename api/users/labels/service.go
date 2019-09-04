package labels

import (
	"context"
	"encoding/base64"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/proto/users/labels/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Service controlls label resource.
type Service struct {
	storage storage.Storage
}

// New creates new labels service instance.
func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// CreateLabel creates new label.
func (s *Service) CreateLabel(ctx context.Context, request *labels.CreateLabelRequest) (*labels.Label, error) {
	if request.Label.Title == "" {
		return nil, status.New(codes.InvalidArgument, "title can't be empty is empty").Err()
	}

	name := resource.ParseName(request.Parent).Child("labels", ksuid.New().String())

	request.Label.Name = name.String()

	raw, err := proto.Marshal(request.Label)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to marshal the label").Err()
	}

	if err := s.storage.Store(name, raw); err != nil {
		return nil, status.New(codes.Internal, "failed to store the label").Err()
	}

	return request.Label, nil
}

// GetLabel returns label by id.
func (s *Service) GetLabel(ctx context.Context, request *labels.GetLabelRequest) (*labels.Label, error) {
	name := resource.ParseName(request.Name)

	raw, err := s.storage.Load(name)
	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, status.New(codes.NotFound, "not found").Err()
	default:
		return nil, status.New(codes.Internal, "failed to load label").Err()
	}

	label := &labels.Label{}
	if err := proto.Unmarshal(raw, label); err != nil {
		return nil, status.New(codes.Internal, "failed to unmarshal the article").Err()
	}

	return label, nil
}

// ListLabels returns list of existing labels.
func (s *Service) ListLabels(ctx context.Context, request *labels.ListLabelsRequest) (*labels.ListLabelsResponse, error) {
	lookFor := resource.ParseName(request.Parent).Child("labels", "*")

	var from *resource.Name
	if request.PageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(request.PageToken)
		if err != nil {
			return nil, status.New(codes.InvalidArgument, "invalid page token").Err()
		}
		from = resource.ParseName(string(decoded))
	}

	dd, err := s.storage.LoadChildren(lookFor, from, int(request.PageSize+1))
	if err != nil {
		return nil, status.New(codes.Internal, "failed to load labels").Err()
	}

	ll := make([]*labels.Label, 0, len(dd))
	for _, d := range dd {
		l := &labels.Label{}
		if err := proto.Unmarshal(d.Data, l); err != nil {
			return nil, status.New(codes.Internal, "failed to unmarshal label").Err()
		}
		ll = append(ll, l)
	}

	var nextPageToken string
	if len(ll) == int(request.PageSize+1) {
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(ll[len(ll)-1].Name))
		ll = ll[:request.PageSize]
	}

	return &labels.ListLabelsResponse{
		Labels:        ll,
		NextPageToken: nextPageToken,
	}, nil
}
