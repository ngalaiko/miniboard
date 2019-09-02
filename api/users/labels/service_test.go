package labels

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/proto/users/labels/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
	"miniboard.app/storage/resource"
)

func TestLabels(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("With a new service", func(t *testing.T) {
		service := New(testDB(ctx, t))
		t.Run("When getting an unknown label", func(t *testing.T) {
			resp, err := service.GetLabel(ctx, &labels.GetLabelRequest{
				Name: resource.NewName("label", "random").String(),
			})

			t.Run("Then error should be not found", func(t *testing.T) {
				assert.Error(t, err)
				assert.Empty(t, resp)

				status, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.NotFound, status.Code())
			})
		})
		t.Run("When creting a label with empty title", func(t *testing.T) {
			resp, err := service.CreateLabel(ctx, &labels.CreateLabelRequest{
				Parent: resource.NewName("users", "name").String(),
				Label:  &labels.Label{},
			})
			t.Run("Invalid argument error should be returned", func(t *testing.T) {
				if assert.Error(t, err) {
					assert.Empty(t, resp)

					status, ok := status.FromError(err)
					assert.True(t, ok)
					assert.Equal(t, codes.InvalidArgument, status.Code())
				}
			})
		})

		t.Run("When creting a label", func(t *testing.T) {
			resp, err := service.CreateLabel(ctx, &labels.CreateLabelRequest{
				Parent: resource.NewName("users", "name").String(),
				Label: &labels.Label{
					Title: "new label",
				},
			})
			t.Run("It should be created", func(t *testing.T) {
				if assert.NoError(t, err) {
					assert.Equal(t, "new label", resp.Title)
					assert.NotEmpty(t, resp.Name)
				}
			})
			t.Run("It should be found", func(t *testing.T) {
				got, err := service.GetLabel(ctx, &labels.GetLabelRequest{
					Name: resp.Name,
				})
				assert.NoError(t, err)
				assert.Equal(t, resp.Name, got.Name)
				assert.Equal(t, resp.Title, got.Title)
			})
		})
	})
}

func testDB(ctx context.Context, t *testing.T) storage.Storage {
	tmpfile, err := ioutil.TempFile("", "bolt")
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}
	go func() {
		<-ctx.Done()
		defer os.Remove(tmpfile.Name()) // clean up
	}()

	db, err := bolt.New(ctx, tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}
	return db
}
