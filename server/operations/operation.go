package operations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/operations/db"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/longrunning"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service manages longrunning api operations.
type Service struct {
	db *db.DB
}

// New returns new service.
func New(sqldb *sql.DB) *Service {
	return &Service{
		db: db.New(sqldb),
	}
}

// Operation is a single long running operation.
type Operation func(context.Context, *longrunning.Operation) (*any.Any, *rpcstatus.Status)

// CreateOperation creates an operation, and runs it.
func (s *Service) CreateOperation(ctx context.Context, metadata *any.Any, operationFunc Operation) (*longrunning.Operation, error) {
	a, _ := actor.FromContext(ctx)
	operation := &longrunning.Operation{
		Name:     fmt.Sprintf("operations/%s", ksuid.New().String()),
		Metadata: metadata,
	}
	if err := s.db.Create(ctx, a.ID, operation); err != nil {
		return nil, fmt.Errorf("failed to create operation: %w", err)
	}

	// todo: wait for graceful shutdown
	go s.runOperation(context.Background(), a.ID, operation, operationFunc)

	return operation, nil
}

func (s *Service) runOperation(ctx context.Context, userID string, operation *longrunning.Operation, operationFunc Operation) {
	result, err := operationFunc(ctx, operation)
	switch err {
	case nil:
		operation.Result = &longrunning.Operation_Response{
			Response: result,
		}
	default:
		operation.Result = &longrunning.Operation_Error{
			Error: err,
		}
	}

	operation.Done = true

	if err := s.db.Update(ctx, userID, operation); err != nil {
		log().Errorf("failed to update operation: %s", err)
	}
}

// GetOperation returns a single operation.
func (s *Service) GetOperation(ctx context.Context, in *longrunning.GetOperationRequest) (*longrunning.Operation, error) {
	a, _ := actor.FromContext(ctx)

	operation, err := s.db.Get(ctx, strings.Replace(in.Name, "operations/", "", -1), a.ID)
	switch {
	case err == nil:
		return operation, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, status.Errorf(codes.NotFound, "not found")
	default:
		return nil, status.Errorf(codes.Internal, "failed to load operation")
	}
}

// ListOperations returns list of operations.
func (s *Service) ListOperations(ctx context.Context, in *longrunning.ListOperationsRequest) (*longrunning.ListOperationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

// DeleteOperation deletes an operation.
func (s *Service) DeleteOperation(ctx context.Context, in *longrunning.DeleteOperationRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

// CancelOperation cancels a running operation.
func (s *Service) CancelOperation(ctx context.Context, in *longrunning.CancelOperationRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

// WaitOperation waits until operations is done, and then returns an operation.
func (s *Service) WaitOperation(ctx context.Context, in *longrunning.WaitOperationRequest) (*longrunning.Operation, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "operations",
	})
}
