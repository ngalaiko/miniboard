package handlers

import (
	"context"

	"github.com/ngalaiko/miniboard/backend/web/sockets"
)

func ItemsLoad(logger logger, itemsService itemsService) sockets.Handler {
	return func(ctx context.Context, req *sockets.Request, respond sockets.Respond, broadcast sockets.Broadcast) {
		refreshed := false
		refreshBeforeRespond := func(resp *sockets.Response) {
			if refreshed {
				respond(resp)
			} else {
				refreshed = true
				resp.Reset = true
				respond(resp)
			}
		}
		ItemsLoadmore(logger, itemsService)(ctx, req, refreshBeforeRespond, broadcast)
	}
}
