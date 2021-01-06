package operations

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_Get_done(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	service := NewService(ctx, &testLogger{}, testDB(ctx, t), nil)

	f := func(ctx context.Context, operation *Operation, status chan<- *Operation) error {
		operation.Done = true
		status <- operation
		return nil
	}

	operation, err := service.Create(ctx, "user", f)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	eventually(t, func() bool {
		doneOperation, err := service.Get(ctx, operation.ID, "user")
		if err != nil {
			t.Fatalf("expected nil, got %s", err)
		}

		return doneOperation.Done && doneOperation.Result == nil
	}, time.Second, 10*time.Millisecond)
}

func Test_Get_panic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	service := NewService(ctx, &testLogger{}, testDB(ctx, t), nil)

	f := func(ctx context.Context, operation *Operation, status chan<- *Operation) error {
		panic("test")
	}

	operation, err := service.Create(ctx, "user", f)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	eventually(t, func() bool {
		doneOperation, err := service.Get(ctx, operation.ID, "user")
		if err != nil {
			t.Fatalf("expected nil, got %s", err)
		}

		return doneOperation.Done &&
			doneOperation.Result != nil &&
			doneOperation.Result.Error.Message == "internal error"
	}, time.Second, 10*time.Millisecond)
}

func Test_Get_error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	service := NewService(ctx, &testLogger{}, testDB(ctx, t), nil)

	f := func(ctx context.Context, operation *Operation, status chan<- *Operation) error {
		return fmt.Errorf("test error")
	}

	operation, err := service.Create(ctx, "user", f)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	eventually(t, func() bool {
		doneOperation, err := service.Get(ctx, operation.ID, "user")
		if err != nil {
			t.Fatalf("expected nil, got %s", err)
		}

		return doneOperation.Done &&
			doneOperation.Result != nil &&
			doneOperation.Result.Error.Message == "test error"
	}, time.Second, 10*time.Millisecond)
}

func Test_Get_no_updates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	service := NewService(ctx, &testLogger{}, testDB(ctx, t), nil)

	f := func(ctx context.Context, operation *Operation, status chan<- *Operation) error {
		return nil
	}

	operation, err := service.Create(ctx, "user", f)
	if err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	eventually(t, func() bool {
		doneOperation, err := service.Get(ctx, operation.ID, "user")
		if err != nil {
			t.Fatalf("expected nil, got %s", err)
		}

		return !doneOperation.Done &&
			doneOperation.Result == nil
	}, time.Second, 10*time.Millisecond)
}

func eventually(t *testing.T, condition func() bool, waitFor time.Duration, tick time.Duration) bool {
	ch := make(chan bool, 1)

	timer := time.NewTimer(waitFor)
	defer timer.Stop()

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for tick := ticker.C; ; {
		select {
		case <-timer.C:
			t.Fatalf("condition never satisfied")
		case <-tick:
			tick = nil
			go func() { ch <- condition() }()
		case v := <-ch:
			if v {
				return true
			}
			tick = ticker.C
		}
	}
}
