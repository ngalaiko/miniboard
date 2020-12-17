package keys

import (
	"bytes"
	"testing"
)

func Test_New(t *testing.T) {
	key, err := New([]byte("der"))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if key.ID == "" {
		t.Fatalf("id is empty")
	}

	if !bytes.Equal(key.PublicDER, []byte("der")) {
		t.Fatalf("expected %s, got %s", []byte("der"), key.PublicDER)
	}
}

func Test_New__empty_key(t *testing.T) {
	_, err := New([]byte(""))
	if err != ErrKeyIsEmpty {
		t.Fatalf("expected %s, got %s", ErrKeyIsEmpty, err)
	}
}
