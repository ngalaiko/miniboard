package resource

import (
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func Test_name(t *testing.T) {
	t.Run("When creating a new name", func(t *testing.T) {
		resourceName := NewName("resource", ksuid.New().String())

		assert.Equal(t, resourceName.Type(), "resource")
		assert.NotEmpty(t, resourceName.ID())

		t.Run("When adding a child", func(t *testing.T) {
			childName := resourceName.Child("child", "id")
			assert.Len(t, childName.Path(), 2)
			assert.Equal(t, childName.Parent(), resourceName)
		})
	})
}

func Test_parseName(t *testing.T) {
	testCases := []struct {
		in  string
		out *Name
	}{
		{
			in: "users/1",
			out: &Name{
				parts: []*idPart{
					{
						ID:   "1",
						Type: "users",
					},
				},
			},
		},
		{
			in: "users/1/accounts/3",
			out: &Name{
				parts: []*idPart{
					{
						ID:   "1",
						Type: "users",
					},
					{
						ID:   "3",
						Type: "accounts",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			got := ParseName(tc.in)
			assert.Equal(t, got, tc.out)
		})
	}
}

func Test_split(t *testing.T) {
	testCases := []struct {
		in    *Name
		first string
		last  string
	}{
		{
			in:    NewName("users", "1"),
			first: "users",
			last:  "1",
		},

		{
			in:    NewName("clients", "2").Child("users", "1"),
			first: "clients/2/users",
			last:  "1",
		},
	}
	for _, tc := range testCases {
		first, last := tc.in.Split()
		assert.Equal(t, first, tc.first)
		assert.Equal(t, last, tc.last)
	}
}
