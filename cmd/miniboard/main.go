package main // import "miniflux.app/cmd"

import (
	"context"

	"miniboard.app/application"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = application.New(ctx)
}
